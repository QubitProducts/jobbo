package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	mesos "github.com/QubitProducts/jobbo/proto/mesos/v1"
	sched "github.com/QubitProducts/jobbo/proto/mesos/v1/scheduler"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
)

func openSchedulerStream(ctx context.Context, masterUrl string, info mesos.FrameworkInfo) (<-chan *sched.Event, string, error) {
	stream, streamId, err := openSchedulerConnection(ctx, masterUrl, info)
	if err != nil {
		return nil, "", errors.Wrap(err, "Subscribe request failed")
	}

	eventChan := make(chan *sched.Event)
	go func() {
		err := readLoop(ctx, stream, eventChan)
		log.Error("Read loop failed:", err)
		close(eventChan)
	}()
	return eventChan, streamId, nil
}

func openSchedulerConnection(ctx context.Context, masterUrl string, info mesos.FrameworkInfo) (io.ReadCloser, string, error) {
	callType := sched.Call_SUBSCRIBE
	subscribeCall := sched.Call{
		Type: &callType,
		Subscribe: &sched.Call_Subscribe{
			FrameworkInfo: &info,
		},
	}
	if info.Id != nil {
		subscribeCall.FrameworkId = info.Id
	}

	data, err := proto.Marshal(&subscribeCall)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to marshal PB")
	}
	buffer := bytes.NewBuffer(data)

	url := fmt.Sprintf("%v/api/v1/scheduler", masterUrl)
	log.V(2).Infof("making subscribe call to %v", url)
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to build request")
	}
	req = req.WithContext(ctx)

	req.Header["Content-Type"] = []string{"application/x-protobuf"}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", errors.Wrap(err, "subscribe request failed")
	}
	log.V(4).Info("Headers:", res)
	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return nil, "", errors.Errorf("invalid status code (%v): %v", res.StatusCode, string(data))
	}

	sIdH, ok := res.Header["Mesos-Stream-Id"]
	if !ok || len(sIdH) != 1 {
		res.Body.Close()
		return nil, "", errors.New("could not find Mesos-Stream-Id header")
	}
	streamId := sIdH[0]
	log.V(3).Infof("streamID: %v", streamId)

	return res.Body, streamId, nil
}

func readLoop(ctx context.Context, resp io.Reader, eventChan chan<- *sched.Event) error {
	log.V(3).Info("starting read loop")

	body := bufio.NewReader(resp)
	for ctx.Err() == nil {
		event, err := readEvent(ctx, body)
		if err != nil {
			return err
		}
		if event != nil {
			eventChan <- event
		}
	}
	return ctx.Err()
}

func readEvent(ctx context.Context, body *bufio.Reader) (*sched.Event, error) {
	log.V(5).Info("Reading header")
	header, err := body.ReadString('\n')
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read header")
	}
	header = strings.Trim(header, "\n")
	// Blank lines?
	if len(header) == 0 {
		return nil, nil
	}
	log.V(5).Info("Header:", header)

	msgLength, err := strconv.Atoi(header)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert header to integer")
	}

	messageBytes := make([]byte, msgLength)
	read := 0
	for read != msgLength {
		n, err := body.Read(messageBytes[read:])
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read message")
		}
		read += n

		log.V(5).Info("Bytes:", read)
		log.V(5).Info("Expected:", msgLength)
		log.V(5).Info("Body:", string(messageBytes))
	}

	var event sched.Event
	err = json.Unmarshal(messageBytes, &event)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal message")
	}
	log.V(4).Info("message:", event)
	log.V(3).Info("recieved ", event.Type)
	return &event, nil
}

func SendCall(ctx context.Context, masterUrl, streamId string, call sched.Call) (*http.Response, error) {
	url := fmt.Sprintf("%v/api/v1/scheduler", masterUrl)

	if log.V(4) {
		data, err := json.Marshal(&call)
		if err != nil {
			log.V(4).Info("Failde to marshal call to json: ", err)
		} else {
			log.V(4).Info("Sending ", string(data))
		}
	}

	data, err := proto.Marshal(&call)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal PB")
	}

	buffer := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build request")
	}
	req = req.WithContext(ctx)

	req.Header["Content-Type"] = []string{"application/x-protobuf"}
	req.Header["Mesos-Stream-Id"] = []string{streamId}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Call request failed")
	}
	defer res.Body.Close()

	log.V(3).Infof("%v %v %v", call.Type, res.StatusCode, url)
	if res.StatusCode != 202 {
		data, _ := ioutil.ReadAll(res.Body)
		return nil, errors.Errorf("unexpected status code (%v): %v", res.StatusCode, string(data))
	}

	return res, nil
}
