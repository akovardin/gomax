package gomax

import (
	"fmt"
	"sync"
	"time"

	"github.com/akovardin/gomax/api"
	apiauth "github.com/akovardin/gomax/api/auth"
	"github.com/akovardin/gomax/api/core"
	apisession "github.com/akovardin/gomax/api/session"
	gomaxauth "github.com/akovardin/gomax/auth"
	"github.com/akovardin/gomax/connection"
	"github.com/akovardin/gomax/dispatch"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/session"
	"github.com/akovardin/gomax/types"
)

type App struct {
	conn       *connection.ConnectionManager
	dispatcher *dispatch.Dispatcher
	api        *api.Facade
	config     *types.ClientConfig
	store      *session.Store
	authFlow   gomaxauth.AuthFlow

	me       *types.Profile
	chats    []*types.Chat
	users    map[types.FlexInt]*types.User
	contacts []*types.User
	messages map[types.FlexInt][]*types.Message

	session *session.SessionInfo

	started  bool
	pingStop chan struct{}
	mu       sync.RWMutex
}

func NewApp(
	conn *connection.ConnectionManager,
	config *types.ClientConfig,
	authFlow gomaxauth.AuthFlow,
	rootRouter *dispatch.Router,
) *App {
	if rootRouter == nil {
		rootRouter = dispatch.NewRouter()
	}

	app := &App{
		conn:     conn,
		config:   config,
		authFlow: authFlow,
		users:    make(map[types.FlexInt]*types.User),
		contacts: make([]*types.User, 0),
		messages: make(map[types.FlexInt][]*types.Message),
		pingStop: make(chan struct{}),
	}

	app.dispatcher = dispatch.NewDispatcher(rootRouter, app)
	app.api = api.NewFacade(app)

	conn.OnEvent = app.onEvent
	conn.OnClose = app.onConnectionLost

	return app
}

func (a *App) Invoke(opcode int, payload map[string]interface{}, cmd int, timeout float64, compress bool) (*protocol.InboundFrame, error) {
	core.LogDebug("invoke opcode=%d cmd=%d timeout=%.0f payload_keys=%v", opcode, cmd, timeout, payloadKeys(payload))
	frame := &protocol.OutboundFrame{
		Ver:     a.conn.Version(),
		Opcode:  opcode,
		Cmd:     cmd,
		Payload: payload,
	}

	dur := time.Duration(timeout * float64(time.Second))
	response, err := a.conn.Request(frame, dur)
	if err != nil {
		return nil, err
	}
	core.LogDebug("response opcode=%d cmd=%d seq=%v", response.Opcode, response.Cmd, response.Seq)
	return response, nil
}

func (a *App) Config() *types.ClientConfig {
	return a.config
}

func (a *App) Token() string {
	if a.session != nil {
		return a.session.Token
	}
	return a.config.Token
}

func (a *App) UpdateToken(oldToken, newToken string) {
	if a.store != nil {
		a.store.UpdateToken(oldToken, newToken)
	}
	if a.session != nil {
		a.session.Token = newToken
	}
}

func (a *App) Me() *types.Profile {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.me
}

func (a *App) Chats() []*types.Chat {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.chats
}

func (a *App) Users() map[types.FlexInt]*types.User {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.users
}

func (a *App) Contacts() []*types.User {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.contacts
}

func (a *App) Messages() map[types.FlexInt][]*types.Message {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.messages
}

func (a *App) CacheChat(chat *types.Chat) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, c := range a.chats {
		if c.ID == chat.ID {
			a.chats[i] = chat
			return
		}
	}
	a.chats = append(a.chats, chat)
}

func (a *App) GetCachedChat(chatID types.FlexInt) *types.Chat {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, c := range a.chats {
		if c.ID == chatID {
			return c
		}
	}
	return nil
}

func (a *App) CacheUser(user *types.User) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.users[user.ID] = user
}

func (a *App) GetCachedUser(userID types.FlexInt) *types.User {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.users[userID]
}

func (a *App) BindMessage(msg *types.Message, chatID types.FlexInt) {
	msg.ChatID = &chatID
}

func (a *App) Start() error {
	if a.started {
		return nil
	}

	a.store = createStore(a.config.WorkDir, a.config.SessionName)

	sessionData, err := a.store.LoadSession()
	if err != nil {
		core.LogWarn("failed to load session: %v", err)
	}

	if sessionData != nil && sessionData.MtInstanceID != "" {
		a.config.Device.MtInstanceID = sessionData.MtInstanceID
	} else if sessionData != nil {
		sessionData.MtInstanceID = a.config.Device.MtInstanceID
	}

	core.LogDebug("opening connection to %s:%d", a.config.Host, a.config.Port)
	if err := a.conn.Open(); err != nil {
		a.store.Close()
		return fmt.Errorf("failed to connect: %w", err)
	}

	handshakeDeviceID := a.config.Device.DeviceID
	if sessionData != nil {
		handshakeDeviceID = sessionData.DeviceID
	}
	core.LogDebug("handshake deviceId=%s", handshakeDeviceID)

	sessionSvc := apisession.NewService(a)
	if err := sessionSvc.Handshake(
		a.config.Device.MtInstanceID,
		a.config.Device.UserAgent,
		handshakeDeviceID,
	); err != nil {
		a.conn.Close()
		a.store.Close()
		return fmt.Errorf("handshake failed: %w", err)
	}

	core.LogDebug("handshake completed")

	go a.pingLoop()

	if sessionData == nil {
		if a.config.Token != "" {
			core.LogDebug("no saved session, using config token")
			sessionData = &session.SessionInfo{
				Token:        a.config.Token,
				DeviceID:     a.config.Device.DeviceID,
				Phone:        a.config.Phone,
				MtInstanceID: a.config.Device.MtInstanceID,
				Sync:         types.DefaultSyncState(),
			}
			if err := a.store.SaveSession(sessionData); err != nil {
				core.LogError("failed to save session: %v", err)
			}
		} else {
			core.LogDebug("no saved session, starting authentication")
			result, err := a.authFlow.Authenticate(a)
			if err != nil {
				a.conn.Close()
				a.store.Close()
				return fmt.Errorf("authentication failed: %w", err)
			}

			if result.Token == "" {
				a.conn.Close()
				a.store.Close()
				return fmt.Errorf("authentication finished without token")
			}

			sessionData = &session.SessionInfo{
				Token:        result.Token,
				DeviceID:     a.config.Device.DeviceID,
				Phone:        a.config.Phone,
				MtInstanceID: a.config.Device.MtInstanceID,
				Sync:         types.DefaultSyncState(),
			}
			if err := a.store.SaveSession(sessionData); err != nil {
				core.LogError("failed to save session: %v", err)
			}
		}
	} else {
		core.LogDebug("loaded saved session deviceId=%s phone=%s", sessionData.DeviceID, sessionData.Phone)
	}

	a.session = sessionData
	core.LogDebug("logging in with token=%s...", sessionData.Token[:min(8, len(sessionData.Token))])

	authSvc := apiauth.NewService(a)
	response, err := authSvc.Login(a.config.Device.UserAgent)
	if err != nil {
		a.conn.Close()
		a.store.Close()
		return fmt.Errorf("login failed: %w", err)
	}

	tokenRotated := false
	if response.Token != nil && *response.Token != a.session.Token {
		a.UpdateToken(a.session.Token, *response.Token)
		tokenRotated = true
	}
	core.LogDebug("login response: token_rotated=%v chats=%d", tokenRotated, len(response.Chats))

	a.me = response.Profile
	a.chats = response.Chats
	if a.me != nil && a.me.Contact != nil {
		a.users[a.me.Contact.ID] = a.me.Contact
	}
	a.contacts = response.Contacts
	for k, v := range response.Messages {
		a.messages[types.FlexInt(k)] = v
	}

	a.updateSyncState(response)

	a.mu.Lock()
	a.started = true
	a.mu.Unlock()

	core.LogInfo("client started profile=%d chats=%d", a.me.Contact.ID, len(a.chats))

	a.dispatcher.OnInternal(dispatch.EventTypeVideoReady)(func(event interface{}, client interface{}) error {
		if sig, ok := event.(*types.VideoUploadSignal); ok {
			id := int(sig.VideoID)
			a.api.Uploads.SignalVideoReady(id)
		}
		return nil
	})
	a.dispatcher.OnInternal(dispatch.EventTypeFileReady)(func(event interface{}, client interface{}) error {
		if sig, ok := event.(*types.FileUploadSignal); ok {
			id := int(sig.FileID)
			a.api.Uploads.SignalFileReady(id)
		}
		return nil
	})

	return nil
}

func (a *App) Close() error {
	a.mu.Lock()
	a.started = false
	a.mu.Unlock()

	if a.pingStop != nil {
		select {
		case <-a.pingStop:
		default:
			close(a.pingStop)
		}
	}

	if a.conn != nil {
		a.conn.Close()
	}

	if a.store != nil {
		a.store.Close()
	}

	return nil
}

func (a *App) pingLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.pingStop:
			return
		case <-ticker.C:
			a.Invoke(
				int(protocol.OpcodePing),
				map[string]interface{}{"interactive": true},
				int(protocol.CommandRequest),
				a.config.RequestTimeout,
				false,
			)
		}
	}
}

func (a *App) onEvent(frame *protocol.InboundFrame) {
	if a.dispatcher != nil {
		a.dispatcher.Dispatch(frame)
	}
}

func (a *App) onConnectionLost(err error) {
	a.mu.Lock()
	a.started = false
	a.mu.Unlock()

	if err != nil {
		core.LogWarn("connection lost: %v", err)
	}
}

func (a *App) IsStarted() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.started
}

func (a *App) updateSyncState(response *types.LoginResponse) {
	if a.session == nil || a.store == nil {
		return
	}
	sync := a.session.Sync
	if response.Time != nil {
		t := *response.Time
		sync.ChatsSync = t
		sync.ContactsSync = t
		sync.DraftsSync = t
		sync.PresenceSync = t
	}
	if response.Config != nil && response.Config.Hash != nil {
		sync.ConfigHash = *response.Config.Hash
	}
	a.session.Sync = sync
	a.store.SaveSession(a.session)
}

func payloadKeys(payload map[string]interface{}) []string {
	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	return keys
}

func createStore(workDir string, dbName string) *session.Store {
	store, err := session.NewStore(workDir, dbName)
	if err != nil {
		core.LogWarn("failed to create store: %v", err)
		return nil
	}
	return store
}
