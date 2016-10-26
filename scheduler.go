package main

import (
	log "github.com/golang/glog"
	"golang.org/x/net/context"

	pb "github.com/QubitProducts/jobbo/proto"
	mesos "github.com/QubitProducts/jobbo/proto/mesos/v1"
	sched "github.com/QubitProducts/jobbo/proto/mesos/v1/scheduler"
	"github.com/pkg/errors"
)

type Task struct {
	TaskId  string
	AgentId string
	State   mesos.TaskState
	UUID    string
}

type Scheduler struct {
	MasterUrl   string
	StreamId    string
	FrameworkId string

	Ledger *JobLedger
	Store  PersistentStore
}

func (s *Scheduler) HandleSubscribed(ctx context.Context, subscription sched.Event_Subscribed) {
	s.FrameworkId = *subscription.FrameworkId.Value
	log.V(1).Infof("subscribed. frameworkID: %v", s.FrameworkId)

	err := s.Store.SetFrameworkID(ctx, s.FrameworkId)
	if err != nil {
		log.Errorf("could not persist framework id to store: %v", err)
	}
}

func (s *Scheduler) HandleUpdate(ctx context.Context, update sched.Event_Update) {
	log.V(2).Infof("got update for %v", *update.Status.TaskId.Value)

	taskName, taskUid, err := decodeTaskID(*update.Status.TaskId.Value)
	if err != nil {
		log.Errorf("could not parse taskID: %v", err)
		return
	}

	status, err := convertMesosStatus(*update.Status.State)
	if err != nil {
		log.Errorf("could not convert task mesos state to jobbo status: %v", err)
		s.ackUpdate(ctx, &update)
		return
	}

	failureReason := ""
	if status == pb.TaskStatus_FAILED {
		failureReason = *update.Status.Message
	}

	err = s.Ledger.MutateEntry(ctx, taskUid, func(s *pb.InstanceState) error {
		s.Status = status
		s.FailureReason = failureReason
		return nil
	})
	if err != nil {
		log.Errorf("could not update ledger state for %v.%v: %v", taskName, taskUid, err)
		return
	}

	s.ackUpdate(ctx, &update)
	log.V(2).Infof("updated %v", *update.Status.TaskId.Value)
}

func (s *Scheduler) HandleOffers(ctx context.Context, offers sched.Event_Offers) {
	if s.FrameworkId == "" {
		log.Info("No framework ID yet")
		return
	}

	log.Infof("Got %v offers!", len(offers.Offers))

	declined := []*mesos.OfferID{}
	for _, o := range offers.Offers {
		job, err := s.FindFittingJob(ctx, o)
		if err != nil {
			log.Errorf("Could not handle offer: %v", err)
		}

		if job == nil {
			declined = append(declined, o.Id)
			continue
		}
		taskInfo := convertJob(job)
		taskInfo.AgentId = o.AgentId

		err = s.launchJob(ctx, o, taskInfo)
		if err != nil {
			log.Errorf("Could not launch job: %v", err)
			break
		}

		err = s.Ledger.StartJob(ctx, job)
		if err != nil {
			log.Errorf("Could not update ledger: %v", err)
			break
		}

		// This might be a bit wierd, but we're only handling one offer at a time. Trust me on this
		break
	}

	log.Info("Declining ", len(declined), " offers")

	err := s.declineOffers(ctx, declined)
	if err != nil {
		log.Error("Failed decline call:", err)
		return
	}
}

func (s *Scheduler) FindFittingJob(ctx context.Context, offer *mesos.Offer) (*pb.Job, error) {
	for _, instance := range s.Ledger.PendingJobSnapshot() {
		if !resourcesFit(offer.Resources, instance.Job.Resources) {
			log.V(2).Infof("Skipping offer %v for %v due to not enough resources", *offer.Id, instance.Job.Metadata.Name)
			continue
		}

		return instance.Job, nil
	}

	return nil, nil
}

func (s *Scheduler) declineOffers(ctx context.Context, offerIDs []*mesos.OfferID) error {
	declineType := sched.Call_DECLINE
	declineCall := sched.Call{
		Type: &declineType,
		FrameworkId: &mesos.FrameworkID{
			Value: &s.FrameworkId,
		},
		Decline: &sched.Call_Decline{
			OfferIds: offerIDs,
		},
	}

	_, err := SendCall(ctx, s.MasterUrl, s.StreamId, declineCall)
	return errors.WithMessage(err, "failed to decline offers")
}

func (s *Scheduler) launchJob(ctx context.Context, offer *mesos.Offer, task *mesos.TaskInfo) error {
	acceptCall := sched.Call{
		Type: sched.Call_ACCEPT.Enum(),
		FrameworkId: &mesos.FrameworkID{
			Value: &s.FrameworkId,
		},
		Accept: &sched.Call_Accept{
			OfferIds: []*mesos.OfferID{offer.Id},
			Operations: []*mesos.Offer_Operation{
				&mesos.Offer_Operation{
					Type: mesos.Offer_Operation_LAUNCH.Enum(),
					Launch: &mesos.Offer_Operation_Launch{
						TaskInfos: []*mesos.TaskInfo{
							task,
						},
					},
				},
			},
		},
	}

	_, err := SendCall(ctx, s.MasterUrl, s.StreamId, acceptCall)
	if err != nil {
		log.Error("Failed to accept call: ", err)
		return err
	}

	return nil
}

func (s *Scheduler) ackUpdate(ctx context.Context, update *sched.Event_Update) error {
	ackCall := sched.Call{
		Type: sched.Call_ACKNOWLEDGE.Enum(),
		FrameworkId: &mesos.FrameworkID{
			Value: &s.FrameworkId,
		},
		Acknowledge: &sched.Call_Acknowledge{
			AgentId: update.Status.AgentId,
			TaskId:  update.Status.TaskId,
			Uuid:    update.Status.Uuid,
		},
	}

	_, err := SendCall(ctx, s.MasterUrl, s.StreamId, ackCall)
	if err != nil {
		log.Errorf("Failed to ack update: %v", err)
		return errors.WithMessage(err, "failed to ack")
	}
	return nil
}

func resourcesFit(offered []*mesos.Resource, required *pb.Resources) bool {
	found := 0
	for _, r := range offered {
		var val float32
		if *r.Type == mesos.Value_SCALAR {
			val = float32(*r.Scalar.Value)
		}

		switch *r.Name {
		case "cpus":
			if val >= required.Cpus {
				found |= 1
			}
		case "mem":
			if val >= required.Mem {
				found |= 2
			}
		case "disk":
			if val >= required.Disk {
				found |= 4
			}
		}
	}
	return found == 7
}
