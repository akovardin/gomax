package dispatch

import (
	"log"
	"sync"

	"github.com/akovardin/gomax/protocol"
)

type Dispatcher struct {
	RootRouter     *Router
	InternalRouter *Router
	Mapper         *EventMapper
	Client         interface{}

	startupTasks []func()
	mu           sync.Mutex
}

func NewDispatcher(rootRouter *Router, app interface{}) *Dispatcher {
	if rootRouter == nil {
		rootRouter = NewRouter()
	}
	return &Dispatcher{
		RootRouter:     rootRouter,
		InternalRouter: NewRouter(),
		Mapper:         NewEventMapper(app),
	}
}

func (d *Dispatcher) BindClient(client interface{}) {
	d.Client = client
}

func (d *Dispatcher) Dispatch(frame *protocol.InboundFrame) {
	eventType := d.resolveEvent(frame)
	var event interface{} = frame

	if eventType != "" && d.Mapper != nil {
		event = d.Mapper.Map(eventType, frame)
	}

	if eventType != "" {
		d.dispatchToRouter(d.InternalRouter, eventType, event)
		d.dispatchToRouter(d.RootRouter, eventType, event)
	}

	d.dispatchToRouter(d.RootRouter, EventTypeRaw, frame)
}

func (d *Dispatcher) resolveEvent(frame *protocol.InboundFrame) EventType {
	if frame.Cmd != int(protocol.CommandRequest) {
		return ""
	}

	if resolver, ok := EventMap[frame.Opcode]; ok {
		return resolver(frame)
	}
	return ""
}

func (d *Dispatcher) dispatchToRouter(router *Router, eventType EventType, event interface{}) {
	if router == nil {
		return
	}

	for _, entry := range router.Handlers[eventType] {
		if d.matches(entry, event) {
			if err := entry.Callback(event, d.Client); err != nil {
				d.handleError(err, eventType, event, router, entry)
			}
		}
	}

	for _, child := range router.Children {
		d.dispatchToRouter(child, eventType, event)
	}
}

func (d *Dispatcher) matches(entry *HandlerEntry, event interface{}) bool {
	for _, filter := range entry.Filters {
		if !filter(event) {
			return false
		}
	}
	return true
}

func (d *Dispatcher) handleError(err error, eventType EventType, event interface{}, originRouter *Router, handler *HandlerEntry) {
	ctx := &ErrorContext{
		Client:    d.Client,
		EventType: eventType,
		Event:     event,
		Handler:   handler,
		Router:    originRouter,
	}

	handled := false
	for _, r := range d.iterRouters() {
		for _, entry := range r.ErrorHandlers {
			if entry.Scope == ErrorScopeLocal && r != originRouter {
				continue
			}
			handled = true
			if err := entry.Callback(err, ctx); err != nil {
				log.Printf("Error in error handler: %v", err)
			}
		}
	}

	if !handled {
		log.Printf("Unhandled error in dispatcher: %v", err)
	}
}

func (d *Dispatcher) iterRouters() []*Router {
	var routers []*Router
	d.collectRouters(d.RootRouter, &routers)
	return routers
}

func (d *Dispatcher) collectRouters(router *Router, result *[]*Router) {
	*result = append(*result, router)
	for _, child := range router.Children {
		d.collectRouters(child, result)
	}
}

func (d *Dispatcher) EmitStart() {
	for _, r := range d.iterRouters() {
		for _, handler := range r.OnStartHandlers {
			if err := handler(d.Client); err != nil {
				log.Printf("Error in start handler: %v", err)
			}
		}
	}
}

func (d *Dispatcher) EmitDisconnect(err error, reconnect bool, delay float64) {
	for _, r := range d.iterRouters() {
		for _, handler := range r.DisconnectHandlers {
			if err := handler(err, reconnect, delay); err != nil {
				log.Printf("Error in disconnect handler: %v", err)
			}
		}
	}
}

func (d *Dispatcher) OnInternal(eventType EventType, filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return d.InternalRouter.On(eventType, filters...)
}
