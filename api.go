package main

import (
	pb "github.com/QubitProducts/jobbo/proto"
	protoTime "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type APIServer struct {
	Ledger    *JobLedger
	masterURL string
}

var _ pb.JobboServer = &APIServer{}

func (a *APIServer) RunJob(ctx context.Context, req *pb.RunJobRequest) (*pb.RunJobReply, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if req == nil || req.Job == nil {
		return nil, errors.Errorf("require .job")
	}

	job := req.Job
	err := validateJob(job)
	if err != nil {
		return nil, errors.WithMessage(err, "could not validate job")
	}

	watermarkJob(job)

	err = a.Ledger.AddPendingJob(ctx, job)
	if err != nil {
		return nil, errors.WithMessage(err, "could not add pending job")
	}

	return &pb.RunJobReply{
		Job: job,
	}, nil
}

func validateJob(job *pb.Job) error {
	if job.Metadata == nil {
		return errors.New("require .job.metadata")
	}
	if job.Metadata.Name == "" {
		return errors.New("require .job.metadata.name")
	}
	if job.Resources == nil {
		return errors.New("require .job.resources")
	}
	if job.Resources.Mem == 0 {
		return errors.New("require .job.resources.mem")
	}
	if job.Resources.Cpus == 0 {
		return errors.New("require .job.resources.cpus")
	}
	if job.Resources.Disk == 0 {
		return errors.New("require .job.resources.disk")
	}
	if job.Command == nil {
		return errors.New("require .job.command")
	}
	if job.Container != nil && job.Container.Image == "" {
		return errors.New("require .job.container.image")
	}

	return nil
}

func watermarkJob(job *pb.Job) {
	job.Metadata.Uid = uuid.NewRandom().String()
	now := time.Now().UTC()
	job.Metadata.CreationTime = &protoTime.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}
}

func (a *APIServer) ListInstances(ctx context.Context, req *pb.ListInstancesRequest) (*pb.ListInstancesReply, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	instances := a.Ledger.InstanceSnapshot(func(e *LedgerEntry) bool {
		return (req.Status == pb.TaskStatus_UNKNOWN || e.state.Status == req.Status) &&
			(req.Name == "" || e.job.Metadata.Name == req.Name) &&
			(req.Uid == "" || e.job.Metadata.Uid == req.Uid)
	})

	return &pb.ListInstancesReply{
		Instances: instances,
	}, nil
}

func (a *APIServer) RemoveInstance(ctx context.Context, req *pb.RemoveInstanceRequest) (*pb.RemoveInstanceReply, error) {
	if req.Status == pb.TaskStatus_SCHEDULED || req.Status == pb.TaskStatus_RUNNING {
		return nil, errors.New("running tasks must stop before they can be removed from jobbo")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	instance, err := a.Ledger.RemoveInstance(ctx, func(e *pb.Instance) bool {
		return e.State.Status == req.Status && e.Job.Metadata.Uid == req.Uid
	})
	if err != nil {
		return nil, errors.WithMessage(err, "could not remove instance")
	}

	return &pb.RemoveInstanceReply{
		Instance: instance,
	}, nil
}

func (a *APIServer) GetInstanceLogs(ctx context.Context, req *pb.GetInstanceLogsRequest) (*pb.GetInstanceLogsReply, error) {
	insts := a.Ledger.InstanceSnapshot(func(e *LedgerEntry) bool {
		return e.job.Metadata.Uid == req.Uid
	})
	if len(insts) == 0 {
		return nil, EntryNotFound
	}

	inst := insts[0]
	if inst.State.Status != pb.TaskStatus_RUNNING && inst.State.Status != pb.TaskStatus_FAILED && inst.State.Status != pb.TaskStatus_FINISHED {
		return nil, errors.New("logs are only availible from tasks that have been scheduled")
	}

	stdout, stderr, err := getInstanceLogs(ctx, a.masterURL, inst)
	if err != nil {
		return nil, errors.WithMessage(err, "could not locate logs")
	}

	return &pb.GetInstanceLogsReply{
		Stderr: stderr,
		Stdout: stdout,
	}, nil
}
