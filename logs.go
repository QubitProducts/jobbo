package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	pb "github.com/QubitProducts/jobbo/proto"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	slavePort = flag.Int("mesos.slave-port", 5051, "Port that slaves expose their HTTP API on")
)

type MasterState struct {
	Slaves []struct {
		ID       string `json:"id"`
		Hostname string `json:"hostname"`
	} `json:"slaves"`
	Frameworks []struct {
		ID    string `json:"id"`
		Tasks []struct {
			ID      string `json:"id"`
			SlaveID string `json:"slave_id"`
		} `json:"tasks"`
		CompletedTasks []struct {
			ID      string `json:"id"`
			SlaveID string `json:"slave_id"`
		} `json:"completed_tasks"`
	} `json:"frameworks"`
}

type SlaveState struct {
	Frameworks          []Framework `json:"frameworks"`
	CompletedFrameworks []Framework `json:"completed_frameworks"`
}
type Framework struct {
	Executors []struct {
		ID        string `json:"id"`
		Directory string `json:"directory"`
	} `json:"executors"`
	CompletedExecutors []struct {
		ID        string `json:"id"`
		Directory string `json:"directory"`
	} `json:"completed_executors"`
}

func getInstanceLogs(ctx context.Context, masterURL string, instance *pb.Instance) (stdout []byte, stderr []byte, err error) {
	slaveHostname, err := getInstanceSlaveHostname(ctx, masterURL, instance)
	if err != nil {
		err = errors.WithMessage(err, "could not find slave")
		return
	}
	slaveURL := fmt.Sprintf("http://%v:%v", slaveHostname, *slavePort)

	dir, err := getInstanceSlaveDirectory(ctx, slaveURL, instance)
	if err != nil {
		err = errors.WithMessage(err, "could not find logs")
		return
	}

	stderr, err = getSlaveFileContents(ctx, slaveURL, fmt.Sprintf("%v/stderr", dir))
	if err != nil {
		err = errors.WithMessage(err, "could not get stderr")
		return
	}
	stdout, err = getSlaveFileContents(ctx, slaveURL, fmt.Sprintf("%v/stdout", dir))
	if err != nil {
		err = errors.WithMessage(err, "could not get stdout")
		return
	}

	return
}

func getInstanceSlaveDirectory(ctx context.Context, slaveURL string, instance *pb.Instance) (string, error) {
	var slave SlaveState
	err := getState(ctx, slaveURL, &slave)
	if err != nil {
		return "", errors.WithMessage(err, "could not find slave state")
	}

	for _, f := range append(slave.Frameworks, slave.CompletedFrameworks...) {
		for _, e := range append(f.Executors, f.CompletedExecutors...) {
			if e.ID == taskID(instance.Job) {
				return e.Directory, nil
			}
		}
	}
	return "", errors.Errorf("Could not find %v on %v", taskID(instance.Job), slaveURL)
}

func getInstanceSlaveHostname(ctx context.Context, masterURL string, instance *pb.Instance) (string, error) {
	instanceID := taskID(instance.Job)

	var master MasterState
	err := getState(ctx, masterURL, &master)
	if err != nil {
		return "", errors.WithMessage(err, "could not get master state")
	}

	slaveID := ""
	for _, f := range master.Frameworks {
		for _, t := range append(f.Tasks, f.CompletedTasks...) {
			if t.ID == instanceID {
				slaveID = t.SlaveID
				break
			}
		}
		if slaveID != "" {
			break
		}
	}
	if slaveID == "" {
		return "", errors.Errorf("could not find instance %v on mesos cluster", instanceID)
	}

	for _, s := range master.Slaves {
		if s.ID == slaveID {
			return s.Hostname, nil
		}
	}
	return "", errors.Errorf("could not find slave %v on mesos cluster", slaveID)
}

func getState(ctx context.Context, baseURL string, target interface{}) error {
	url := fmt.Sprintf("%v/state.json", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, "could not create request to %v", url)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "could not perform request to %v", url)
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("invalid status code on call to %v (%v): %v", url, resp.StatusCode, string(data))
	}
	glog.V(4).Infof("GET %v %v", resp.StatusCode, url)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read response")
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response")
	}

	return nil
}

func getSlaveFileContents(ctx context.Context, slaveURL string, filepath string) ([]byte, error) {
	url := fmt.Sprintf("%v/files/download?path=%v", slaveURL, url.QueryEscape(filepath))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create request to %v", url)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not perform request to %v", url)
	}
	defer resp.Body.Close()
	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("invalid status code on call to %v (%v): %v", url, resp.StatusCode, string(data))
	}
	glog.V(4).Infof("GET %v %v", resp.StatusCode, url)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response")
	}

	return data, nil
}
