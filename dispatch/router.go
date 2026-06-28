package dispatch

type EventType = string

type FilterCallback func(event interface{}) bool

type HandlerCallback func(event interface{}, client interface{}) error

type StartCallback func(client interface{}) error

type ErrorCallback func(err error, ctx *ErrorContext) error

type DisconnectCallback func(err error, reconnect bool, delay float64) error

type ErrorScope string

const (
	ErrorScopeGlobal ErrorScope = "GLOBAL"
	ErrorScopeLocal  ErrorScope = "LOCAL"
)

type ErrorContext struct {
	Client    interface{}
	EventType EventType
	Event     interface{}
	Handler   interface{}
	Router    *Router
}

type HandlerEntry struct {
	Callback HandlerCallback
	Filters  []FilterCallback
}

type ErrorEntry struct {
	Callback ErrorCallback
	Scope    ErrorScope
}

type Router struct {
	Handlers           map[EventType][]*HandlerEntry
	Children           []*Router
	OnStartHandlers    []StartCallback
	ErrorHandlers      []*ErrorEntry
	DisconnectHandlers []DisconnectCallback
}

func NewRouter() *Router {
	return &Router{
		Handlers: make(map[EventType][]*HandlerEntry),
	}
}

func (r *Router) On(eventType EventType, filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return func(callback HandlerCallback) HandlerCallback {
		r.Handlers[eventType] = append(r.Handlers[eventType], &HandlerEntry{
			Callback: callback,
			Filters:  filters,
		})
		return callback
	}
}

func (r *Router) IncludeRouter(router *Router) {
	r.Children = append(r.Children, router)
}

func (r *Router) OnStart() func(StartCallback) StartCallback {
	return func(callback StartCallback) StartCallback {
		r.OnStartHandlers = append(r.OnStartHandlers, callback)
		return callback
	}
}

func (r *Router) OnError(scope ErrorScope) func(ErrorCallback) ErrorCallback {
	return func(callback ErrorCallback) ErrorCallback {
		r.ErrorHandlers = append(r.ErrorHandlers, &ErrorEntry{
			Callback: callback,
			Scope:    scope,
		})
		return callback
	}
}

func (r *Router) OnDisconnect() func(DisconnectCallback) DisconnectCallback {
	return func(callback DisconnectCallback) DisconnectCallback {
		r.DisconnectHandlers = append(r.DisconnectHandlers, callback)
		return callback
	}
}

func (r *Router) OnMessage(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeMessageNew, filters...)
}
func (r *Router) OnMessageEdit(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeMessageEdit, filters...)
}
func (r *Router) OnMessageDelete(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeMessageDelete, filters...)
}
func (r *Router) OnMessageRead(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeMessageRead, filters...)
}
func (r *Router) OnTyping(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeTyping, filters...)
}
func (r *Router) OnPresence(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypePresence, filters...)
}
func (r *Router) OnReactionUpdate(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeReactionUpdate, filters...)
}
func (r *Router) OnChatUpdate(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeChatUpdate, filters...)
}
func (r *Router) OnRaw(filters ...FilterCallback) func(HandlerCallback) HandlerCallback {
	return r.On(EventTypeRaw, filters...)
}

const (
	EventTypeMessageNew     EventType = "MESSAGE_NEW"
	EventTypeMessageEdit    EventType = "MESSAGE_EDIT"
	EventTypeMessageDelete  EventType = "MESSAGE_DELETE"
	EventTypeMessageRead    EventType = "MESSAGE_READ"
	EventTypeTyping         EventType = "TYPING"
	EventTypePresence       EventType = "PRESENCE"
	EventTypeReactionUpdate EventType = "REACTION_UPDATE"
	EventTypeChatUpdate     EventType = "CHAT_UPDATE"
	EventTypeUserUpdate     EventType = "USER_UPDATE"
	EventTypeVideoReady     EventType = "VIDEO_READY"
	EventTypeFileReady      EventType = "FILE_READY"
	EventTypeRaw            EventType = "RAW"
	EventTypeOnStart        EventType = "ON_START"
)
