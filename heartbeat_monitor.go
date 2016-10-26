package main

import (
	"time"

	log "github.com/golang/glog"
	"golang.org/x/net/context"

	sched "github.com/QubitProducts/jobbo/proto/mesos/v1/scheduler"
)

const (
	DEFAULT_HEARTBEAT_INTERVAL = 30 * time.Second
)

type HeartbeatMonitor struct {
	CancelFunc    func()
	GracePeriod   time.Duration
	lastHeartbeat time.Time
	interval      time.Duration
}

func (h *HeartbeatMonitor) HandleSubscribed(ctx context.Context, subscription sched.Event_Subscribed) {
	if subscription.HeartbeatIntervalSeconds != nil {
		h.interval = time.Duration((*subscription.HeartbeatIntervalSeconds) * float64(time.Second))
	} else {
		log.V(2).Infof("assuming default heartbeat interval of %v", DEFAULT_HEARTBEAT_INTERVAL)
		h.interval = DEFAULT_HEARTBEAT_INTERVAL
	}
	log.V(4).Info("Heartbeat every ", h.interval)

	go func() {
		for range time.Tick(h.interval / 2) {
			if ctx.Err() != nil {
				log.V(2).Info("Stopping heartbeat monitor loop")
				break
			}

			due := h.lastHeartbeat.Add(h.interval).Add(h.GracePeriod)
			nowish := time.Now()
			if nowish.After(due) {
				log.V(1).Infof("canceling due to missed heartbeat (%v after last heartbeat deadline)", nowish.Sub(due))

				h.CancelFunc()
			}
		}
	}()
}

func (h *HeartbeatMonitor) HandleHeartbeat(ctx context.Context) {
	log.V(4).Info("Updating last heartbeat time")
	h.lastHeartbeat = time.Now()
}
