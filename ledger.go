package main

import (
	"context"
	pb "github.com/QubitProducts/jobbo/proto"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"sync"
)

const (
	PENDING TaskState = iota
	SCHEDULED
	RUNNING
	FAILED
	FINISHED
)

var (
	EntryNotFound = errors.New("could not find ledger entry")
)

type TaskState byte

type JobLedger struct {
	lock *sync.RWMutex

	entries []*LedgerEntry
	store   PersistentStore
}

// Basically an `pb.Instance` with a lock
type LedgerEntry struct {
	lock *sync.RWMutex

	job   *pb.Job
	state *pb.InstanceState
}

func (le *LedgerEntry) Instance() *pb.Instance {
	return &pb.Instance{
		Job:   le.job,
		State: le.state,
	}
}

func NewJobLedger(ctx context.Context, store PersistentStore) (*JobLedger, error) {
	l := &JobLedger{
		lock:    &sync.RWMutex{},
		entries: []*LedgerEntry{},
		store:   store,
	}

	err := l.Sync(ctx)
	return l, errors.WithMessage(err, "could not sync ledger")
}

func (l *JobLedger) Sync(ctx context.Context) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	instances, err := l.store.ListInstances(ctx)
	if err != nil {
		return errors.WithMessage(err, "could not retrieve instances from persistent store")
	}

	entries := []*LedgerEntry{}
	for _, i := range instances {
		entries = append(entries, &LedgerEntry{
			lock:  &sync.RWMutex{},
			job:   i.Job,
			state: i.State,
		})
	}

	l.entries = entries
	return nil
}

func (l *JobLedger) MutateEntry(ctx context.Context, uid string, mut func(*pb.InstanceState) error) error {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, e := range l.entries {
		e.lock.RLock()

		// We haven't deferred unlock here, so kinda important we don't panic
		// Safe to unlock and relock after this, due to UID being immutable
		if e != nil && e.job != nil && e.job.Metadata != nil && e.job.Metadata.Uid != uid {
			e.lock.RUnlock()
			continue
		}

		e.lock.RUnlock()
		e.lock.Lock()
		defer e.lock.Unlock()

		stateCp := proto.Clone(e.state).(*pb.InstanceState)

		err := mut(stateCp)
		if err != nil {
			return errors.Wrap(err, "instance state mutation failed")
		}
		inst := e.Instance()
		inst.State = stateCp

		err = l.store.SetInstance(ctx, inst)
		if err != nil {
			return errors.Wrap(err, "could not update persistent store")
		}
		e.state = stateCp

		return nil
	}
	return EntryNotFound
}

func (l *JobLedger) RemoveInstance(ctx context.Context, pred func(*pb.Instance) bool) (*pb.Instance, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	ix := -1
	for i, ent := range l.entries {
		if pred(ent.Instance()) {
			ix = i
			break
		}
	}
	if ix == -1 {
		return nil, EntryNotFound
	}

	inst := l.entries[ix].Instance()
	err := l.store.RemoveInstance(ctx, inst)
	if err != nil {
		return nil, errors.Wrap(err, "could not persist removal to store")
	}

	l.entries = append(l.entries[0:ix], l.entries[ix+1:]...)

	return inst, nil
}

func (l *JobLedger) InstanceSnapshot(pred func(*LedgerEntry) bool) []*pb.Instance {
	l.lock.RLock()
	defer l.lock.RUnlock()
	instances := make([]*pb.Instance, 0, len(l.entries))

	for _, e := range l.entries {
		func() {
			e.lock.RLock()
			defer e.lock.RUnlock()

			if pred(e) {
				instances = append(instances, &pb.Instance{
					Job:   proto.Clone(e.job).(*pb.Job),
					State: proto.Clone(e.state).(*pb.InstanceState),
				})
			}
		}()
	}
	return instances
}

func (l *JobLedger) PendingJobSnapshot() []*pb.Instance {
	return l.InstanceSnapshot(func(e *LedgerEntry) bool {
		return e.state.Status == pb.TaskStatus_PENDING
	})
}

func (l *JobLedger) AddPendingJob(ctx context.Context, job *pb.Job) error {
	glog.V(2).Infof("Adding %v-%v to pending jobs", job.Metadata.Name, job.Metadata.Uid)
	l.lock.Lock()
	defer l.lock.Unlock()

	entry := &LedgerEntry{
		lock: &sync.RWMutex{},
		job:  proto.Clone(job).(*pb.Job),
		state: &pb.InstanceState{
			Status: pb.TaskStatus_PENDING,
		},
	}

	err := l.store.SetInstance(ctx, entry.Instance())
	if err != nil {
		return errors.Wrap(err, "failed to persist to store")
	}

	l.entries = append(l.entries, entry)
	glog.V(2).Infof("Added %v-%v to pending jobs", job.Metadata.Name, job.Metadata.Uid)

	return nil
}

func (l *JobLedger) StartJob(ctx context.Context, job *pb.Job) error {
	glog.V(2).Infof("Removing %v-%v from pending jobs", job.Metadata.Name, job.Metadata.Uid)
	glog.V(2).Infof("Adding %v-%v to running jobs", job.Metadata.Name, job.Metadata.Uid)

	err := l.MutateEntry(ctx, job.Metadata.Uid, func(s *pb.InstanceState) error {
		s.Status = pb.TaskStatus_SCHEDULED
		return nil
	})
	return errors.WithMessage(err, "could not change job "+job.Metadata.Uid+" state to SCHEDULED")
}

func (l *JobLedger) SetInstanceStatus(ctx context.Context, uid string, status pb.TaskStatus) error {
	err := l.MutateEntry(ctx, uid, func(s *pb.InstanceState) error {
		s.Status = status
		return nil
	})
	return errors.WithMessage(err, "could not update job "+uid+" to status "+string(status))
}
