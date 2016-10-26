package main

import (
	"flag"
	"net"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	pb "github.com/QubitProducts/jobbo/proto"
	mesos "github.com/QubitProducts/jobbo/proto/mesos/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

func init() {
	flag.Parse()
}

const (
	CPUS_PER_EXECUTOR = 0.01
	CPUS_PER_TASK     = 0.1
	MEM_PER_EXECUTOR  = 64
	MEM_PER_TASK      = 64
	PORTS_PER_TASK    = 1

	MESOS_STREAM_ID_CTX_KEY = "mesos-stream-id-ctxkey"
)

var (
	grpcAddress = flag.String("grpc.addr", ":7281", "Bind address for the GRPC server")

	user            = flag.String("mesos.user", "etcd", "Mesos user to run as")
	frameworkName   = flag.String("mesos.framework", "etcd", "Mesos framework to register as")
	masterUrl       = flag.String("mesos.master", "", "HTTP URL of Mesos master")
	failoverTimeout = flag.Duration("mesos.failover-timeout", time.Hour*24*7, "Timeout for mesos framework reconnection")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			glog.Errorf("debug server could not listen on :8080: %v", err)
		}
		cancel()
	}()

	store, err := newEtcdStore(ctx)
	if err != nil {
		glog.Errorf("could not connect to etcd: %v", err)
		os.Exit(1)
	}
	ledger, err := NewJobLedger(ctx, store)
	if err != nil {
		glog.Errorf("could not create job ledger: %v", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer cancel()

		err := runMesosScheduler(ctx, store, ledger)
		if err != nil && errors.Cause(err) != context.Canceled {
			glog.Errorf("Mesos scheduler failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		defer cancel()

		err := runGRPCServer(ctx, ledger, *grpcAddress)
		if err != nil && errors.Cause(err) != context.Canceled {
			glog.Errorf("GRPC server failed: %v", err)
		}
	}()

	wg.Wait()
}

func runGRPCServer(ctx context.Context, ledger *JobLedger, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %v", address)
	}
	s := grpc.NewServer()
	pb.RegisterJobboServer(s, &APIServer{ledger, *masterUrl})

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()
	err = s.Serve(lis)
	if err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}

func runMesosScheduler(ctx context.Context, store PersistentStore, ledger *JobLedger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var frameworkID *mesos.FrameworkID
	id, err := store.GetFrameworkID(ctx)
	if err != nil {
		return errors.WithMessage(err, "cannot get framework id from persistent store")
	}
	if id != "" {
		glog.V(1).Infof("found existing framework id: %v", id)
		frameworkID = &mesos.FrameworkID{
			Value: &id,
		}
	}

	timeout := (*failoverTimeout).Seconds()
	eventStream, streamId, err := openSchedulerStream(ctx, *masterUrl, mesos.FrameworkInfo{
		User:            user,
		Name:            frameworkName,
		Id:              frameworkID,
		FailoverTimeout: &timeout,
	})
	if err != nil {
		return errors.WithMessage(err, "scheduler stream initialisation failed")
	}

	muxer := &EventMuxer{}

	heartbeatMonitor := &HeartbeatMonitor{CancelFunc: cancel, GracePeriod: time.Second * 10}
	muxer.Handle(heartbeatMonitor)

	s := &Scheduler{StreamId: streamId, MasterUrl: *masterUrl, Ledger: ledger, Store: store}
	muxer.Handle(s)

	err = muxer.Run(ctx, eventStream, streamId)
	if err != nil {
		return errors.WithMessage(err, "event muxer exited")
	}

	return nil
}
