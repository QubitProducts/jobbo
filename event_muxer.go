package main

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	sched "github.com/QubitProducts/jobbo/proto/mesos/v1/scheduler"
)

type EventMuxer struct {
	SubscribedHandlers []func(context.Context, sched.Event_Subscribed)
	OffersHandlers     []func(context.Context, sched.Event_Offers)
	RescindHandlers    []func(context.Context, sched.Event_Rescind)
	UpdateHandlers     []func(context.Context, sched.Event_Update)
	MessageHandlers    []func(context.Context, sched.Event_Message)
	FailureHandlers    []func(context.Context, sched.Event_Failure)
	ErrorHandlers      []func(context.Context, sched.Event_Error)
	HeartbeatHandlers  []func(context.Context)
}

func (m *EventMuxer) Run(ctx context.Context, eventStream <-chan *sched.Event, streamId string) error {
	handlerCtx := context.WithValue(ctx, MESOS_STREAM_ID_CTX_KEY, streamId)

	for event := range eventStream {
		if ctx.Err() != nil {
			break
		}

		switch *event.Type {
		case sched.Event_SUBSCRIBED:
			for _, h := range m.SubscribedHandlers {
				h(handlerCtx, *event.Subscribed)
			}
		case sched.Event_OFFERS:
			for _, h := range m.OffersHandlers {
				go h(handlerCtx, *event.Offers)
			}
		case sched.Event_RESCIND:
			for _, h := range m.RescindHandlers {
				go h(handlerCtx, *event.Rescind)
			}
		case sched.Event_UPDATE:
			for _, h := range m.UpdateHandlers {
				go h(handlerCtx, *event.Update)
			}
		case sched.Event_MESSAGE:
			for _, h := range m.MessageHandlers {
				go h(handlerCtx, *event.Message)
			}
		case sched.Event_FAILURE:
			for _, h := range m.FailureHandlers {
				go h(handlerCtx, *event.Failure)
			}
		case sched.Event_ERROR:
			for _, h := range m.ErrorHandlers {
				go h(handlerCtx, *event.Error)
			}
		case sched.Event_HEARTBEAT:
			for _, h := range m.HeartbeatHandlers {
				go h(handlerCtx)
			}
		default:
			return errors.New("Unknown event type from Mesos")
		}
	}
	return nil
}

func (m *EventMuxer) HandleSubscribed(f func(context.Context, sched.Event_Subscribed)) {
	m.SubscribedHandlers = append(m.SubscribedHandlers, f)
}
func (m *EventMuxer) HandleOffers(f func(context.Context, sched.Event_Offers)) {
	m.OffersHandlers = append(m.OffersHandlers, f)
}
func (m *EventMuxer) HandleRescind(f func(context.Context, sched.Event_Rescind)) {
	m.RescindHandlers = append(m.RescindHandlers, f)
}
func (m *EventMuxer) HandleUpdate(f func(context.Context, sched.Event_Update)) {
	m.UpdateHandlers = append(m.UpdateHandlers, f)
}
func (m *EventMuxer) HandleMessage(f func(context.Context, sched.Event_Message)) {
	m.MessageHandlers = append(m.MessageHandlers, f)
}
func (m *EventMuxer) HandleFailure(f func(context.Context, sched.Event_Failure)) {
	m.FailureHandlers = append(m.FailureHandlers, f)
}
func (m *EventMuxer) HandleError(f func(context.Context, sched.Event_Error)) {
	m.ErrorHandlers = append(m.ErrorHandlers, f)
}
func (m *EventMuxer) HandleHeartbeat(f func(context.Context)) {
	m.HeartbeatHandlers = append(m.HeartbeatHandlers, f)
}

type SubscribedHandler interface {
	HandleSubscribed(context.Context, sched.Event_Subscribed)
}
type OffersHandler interface {
	HandleOffers(context.Context, sched.Event_Offers)
}
type RescindHandler interface {
	HandleRescind(context.Context, sched.Event_Rescind)
}
type UpdateHandler interface {
	HandleUpdate(context.Context, sched.Event_Update)
}
type MessageHandler interface {
	HandleMessage(context.Context, sched.Event_Message)
}
type FailureHandler interface {
	HandleFailure(context.Context, sched.Event_Failure)
}
type ErrorHandler interface {
	HandleError(context.Context, sched.Event_Error)
}
type HeartbeatHandler interface {
	HandleHeartbeat(context.Context)
}

func (m *EventMuxer) Handle(h interface{}) {
	if th, ok := h.(SubscribedHandler); ok {
		m.HandleSubscribed(th.HandleSubscribed)
	}
	if th, ok := h.(OffersHandler); ok {
		m.HandleOffers(th.HandleOffers)
	}
	if th, ok := h.(RescindHandler); ok {
		m.HandleRescind(th.HandleRescind)
	}
	if th, ok := h.(UpdateHandler); ok {
		m.HandleUpdate(th.HandleUpdate)
	}
	if th, ok := h.(MessageHandler); ok {
		m.HandleMessage(th.HandleMessage)
	}
	if th, ok := h.(FailureHandler); ok {
		m.HandleFailure(th.HandleFailure)
	}
	if th, ok := h.(ErrorHandler); ok {
		m.HandleError(th.HandleError)
	}
	if th, ok := h.(HeartbeatHandler); ok {
		m.HandleHeartbeat(th.HandleHeartbeat)
	}
}
