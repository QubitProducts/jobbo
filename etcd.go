package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	pb "github.com/QubitProducts/jobbo/proto"
	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	etcdDiscoverySRV       = flag.String("etcd.discovery.srv", "", "SRV record to use to discover ETCD")
	etcdDiscoveryEndpoints = flag.String("etcd.discovery.endpoints", "", "Endpoint list to use to discover ETCD")
	etcdAuthUsername       = flag.String("etcd.auth.username", "", "Username used to auth with ETCD")
	etcdAuthPassword       = flag.String("etcd.auth.password", "", "Password used to auth with ETCD")
	etcdAuthCA             = flag.String("etcd.auth.ca", "", "CA certificate used to connect to ETCD")
	etcdAuthCert           = flag.String("etcd.auth.cert", "", "Client certificate used to connect to ETCD")
	etcdAuthKey            = flag.String("etcd.auth.key", "", "Client key used to connect to ETCD")
	etcdPrefix             = flag.String("etcd.prefix", "/jobbo", "Prefix to use when calculating ETCD keys")
)

type etcdStore struct {
	client etcd.Client
	prefix string
}

func newEtcdStore(ctx context.Context) (*etcdStore, error) {
	conf, err := etcdConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "could not load etcd config")
	}

	client, err := etcd.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "could not create etcd client")
	}

	err = client.Sync(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not perform initial etcd endpoint sync")
	}

	go func() {
		for ctx.Err() == nil {
			func() {
				ctx, cancel := context.WithTimeout(ctx, time.Second*30)
				defer cancel()

				err := client.Sync(ctx)
				if err != nil {
					glog.Errorf("failed to sync etcd endpoints: %v", err)
				} else {
					glog.V(4).Infof("synced etcd endpoint list")
				}
			}()

			select {
			case <-time.After(time.Second * 60):
			case <-ctx.Done():
			}
		}
	}()

	return &etcdStore{
		client: client,
		prefix: *etcdPrefix,
	}, nil
}

func etcdConfig() (etcd.Config, error) {
	https := false
	certificates := []tls.Certificate{}
	if *etcdAuthCert != "" {
		if *etcdAuthKey == "" {
			return etcd.Config{}, errors.New("-etcd.auth.cert requires -etcd.auth.key")
		}
		https = true

		cert, err := tls.LoadX509KeyPair(*etcdAuthCert, *etcdAuthKey)
		if err != nil {
			return etcd.Config{}, errors.Wrap(err, "could not load etcd client key pair")
		}

		certificates = append(certificates, cert)
	}

	caCertPool := x509.NewCertPool()
	if *etcdAuthCA != "" {
		https = true
		caCert, err := ioutil.ReadFile(*etcdAuthCA)
		if err != nil {
			return etcd.Config{}, errors.Wrap(err, "could not load etcd CA")
		}

		caCertPool.AppendCertsFromPEM(caCert)
	}

	endpoints := []string{}
	if *etcdDiscoverySRV != "" {
		discoverer := etcd.NewSRVDiscover()
		var err error
		endpoints, err = discoverer.Discover(*etcdDiscoverySRV)
		if err != nil {
			return etcd.Config{}, errors.Wrapf(err, "could not resolve discovery SRV %v", etcdDiscoverySRV)
		}
	} else if *etcdDiscoveryEndpoints != "" {
		endpoints = strings.Split(*etcdDiscoveryEndpoints, ",")
	} else {
		return etcd.Config{}, errors.New("No etcd discovery mechanism specified (-etcd.discovery.srv, -etcd.discovery.endpoints)")
	}

	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return etcd.Config{}, errors.Wrapf(err, "could not parse initial node %v", e)
		}
		if https && u.Scheme != "https" {
			return etcd.Config{}, errors.Errorf("HTTPS configured, but non https address found in initial nodes: %v", e)
		}
	}

	transport := etcd.DefaultTransport
	if https {
		tlsConfig := &tls.Config{
			Certificates: certificates,
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return etcd.Config{
		Endpoints:     endpoints,
		Transport:     transport,
		Username:      *etcdAuthUsername,
		Password:      *etcdAuthPassword,
		SelectionMode: etcd.EndpointSelectionPrioritizeLeader,
	}, nil
}

func (e *etcdStore) GetFrameworkID(ctx context.Context) (string, error) {
	kAPI := etcd.NewKeysAPI(e.client)
	path := fmt.Sprintf("%v/framework-id", e.prefix)

	resp, err := kAPI.Get(ctx, path, nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return "", nil
		} else {
			return "", errors.Wrap(err, "could not get framework id")
		}
	}

	if resp.Node.Dir {
		return "", errors.Errorf("%v is a directory", path)
	}

	return resp.Node.Value, nil
}

func (e *etcdStore) SetFrameworkID(ctx context.Context, id string) error {
	kAPI := etcd.NewKeysAPI(e.client)
	path := fmt.Sprintf("%v/framework-id", e.prefix)

	_, err := kAPI.Set(ctx, path, id, nil)
	if err != nil {
		return errors.Wrap(err, "could not set framework id")
	}

	return nil
}

func (e *etcdStore) SetInstance(ctx context.Context, inst *pb.Instance) error {
	data, err := proto.Marshal(inst)
	if err != nil {
		return errors.Wrap(err, "could not marshal job")
	}

	kAPI := etcd.NewKeysAPI(e.client)
	path := fmt.Sprintf("%v/instances/%v", e.prefix, inst.Job.Metadata.Uid)
	glog.V(2).Infof("persisting %v to %v", taskID(inst.Job), path)

	_, err = kAPI.Set(ctx, path, base64.StdEncoding.EncodeToString(data), nil)
	if err != nil {
		return errors.Wrapf(err, "could not set %v", path)
	}

	return nil
}

func (e *etcdStore) ListInstances(ctx context.Context) ([]*pb.Instance, error) {
	kAPI := etcd.NewKeysAPI(e.client)
	path := fmt.Sprintf("%v/instances/", e.prefix)

	resp, err := kAPI.Get(ctx, path, nil)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			return []*pb.Instance{}, nil
		} else {
			return []*pb.Instance{}, errors.Wrap(err, "could not set framework id")
		}
	}

	if !resp.Node.Dir {
		return []*pb.Instance{}, errors.Errorf("%v is not a directory", path)
	}

	insts := []*pb.Instance{}
	for _, n := range resp.Node.Nodes {
		data, err := base64.StdEncoding.DecodeString(n.Value)
		if err != nil {
			return insts, errors.Wrapf(err, "could not decode instance at %v", n.Key)
		}

		var inst pb.Instance
		err = proto.Unmarshal(data, &inst)
		if err != nil {
			return insts, errors.Wrapf(err, "could not unmarshal instance at %v", n.Key)
		}
		insts = append(insts, &inst)
	}

	return insts, nil
}

func (e *etcdStore) RemoveInstance(ctx context.Context, inst *pb.Instance) error {
	kAPI := etcd.NewKeysAPI(e.client)
	path := fmt.Sprintf("%v/instances/%v", e.prefix, inst.Job.Metadata.Uid)
	glog.V(2).Infof("removing %v from %v", taskID(inst.Job), path)

	_, err := kAPI.Delete(ctx, path, nil)
	if err != nil {
		return errors.Wrapf(err, "could not delete %v", path)
	}

	return nil
}
