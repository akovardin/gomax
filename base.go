package gomax

import (
	"fmt"
	"time"

	"github.com/akovardin/gomax/api"
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/auth"
	"github.com/akovardin/gomax/connection"
	"github.com/akovardin/gomax/dispatch"
	"github.com/akovardin/gomax/types"
)

type BaseClient struct {
	ExtraConfig *ExtraConfig
	Phone       string
	SessionName string
	WorkDir     string

	config     *types.ClientConfig
	connection *connection.ConnectionManager
	app        *App
	authFlow   auth.AuthFlow
	router     *dispatch.Router
	buildConn  func() *connection.ConnectionManager
}

func (c *BaseClient) Me() *types.Profile {
	if c.app == nil {
		return nil
	}
	return c.app.Me()
}

func (c *BaseClient) Chats() []*types.Chat {
	if c.app == nil {
		return nil
	}
	return c.app.Chats()
}

func (c *BaseClient) Contacts() []*types.User {
	if c.app == nil {
		return nil
	}
	return c.app.Contacts()
}

func (c *BaseClient) Messages() map[types.FlexInt][]*types.Message {
	if c.app == nil {
		return nil
	}
	return c.app.Messages()
}

func (c *BaseClient) API() *api.Facade {
	if c.app == nil {
		return nil
	}
	return c.app.api
}

func (c *BaseClient) BuildConfig(phone string, userAgent *types.MobileUserAgentPayload) *types.ClientConfig {
	if c.ExtraConfig == nil {
		c.ExtraConfig = DefaultExtraConfig()
	}

	deviceID := c.ExtraConfig.DeviceID
	if deviceID == "" {
		deviceID = generateUUID()
	}

	mtInstanceID := c.ExtraConfig.MtInstanceID
	if mtInstanceID == "" {
		mtInstanceID = generateUUID()
	}

	return &types.ClientConfig{
		Phone:              phone,
		SessionName:        c.SessionName,
		WorkDir:            c.WorkDir,
		Token:              c.ExtraConfig.Token,
		Host:               c.ExtraConfig.Host,
		Port:               c.ExtraConfig.Port,
		UseSSL:             c.ExtraConfig.UseSSL,
		ProtocolVersion:    10,
		RequestTimeout:     c.ExtraConfig.RequestTimeout,
		LogLevel:           c.ExtraConfig.LogLevel,
		Telemetry:          c.ExtraConfig.Telemetry,
		Sync:               c.ExtraConfig.Sync,
		Proxy:              c.ExtraConfig.Proxy,
		RegistrationConfig: c.ExtraConfig.RegistrationConfig,
		Device: &types.DeviceConfig{
			MtInstanceID:    mtInstanceID,
			DeviceID:        deviceID,
			UserAgent:       userAgent,
			ClientSessionID: 1,
		},
	}
}

func (c *BaseClient) InitRuntime(authFlow auth.AuthFlow) {
	c.authFlow = authFlow
	c.connection = c.buildConn()
	c.router = dispatch.NewRouter()
	c.app = NewApp(c.connection, c.config, c.authFlow, c.router)
	c.app.dispatcher.BindClient(c)
}

func (c *BaseClient) ResetRuntime() {
	c.connection = c.buildConn()
	c.app = NewApp(c.connection, c.config, c.authFlow, c.router)
	c.app.dispatcher.BindClient(c)
}

func (c *BaseClient) Start() error {
	core.SetLogLevel(c.config.LogLevel)
	for {
		err := c.app.Start()
		if err != nil {
			c.Close()
			return err
		}

		if !c.app.IsStarted() {
			c.Close()
			return nil
		}

		c.app.dispatcher.EmitStart()

		if err := c.connection.WaitClosed(); err != nil {
			c.Close()

			c.app.dispatcher.EmitDisconnect(err, c.ExtraConfig.Reconnect, c.ExtraConfig.ReconnectDelay)

			if !c.ExtraConfig.Reconnect {
				return err
			}

			time.Sleep(time.Duration(c.ExtraConfig.ReconnectDelay * float64(time.Second)))
			c.ResetRuntime()
			continue
		}

		c.Close()
		return nil
	}
}

func (c *BaseClient) Close() error {
	if c.app != nil {
		return c.app.Close()
	}
	return nil
}

func (c *BaseClient) Stop() error {
	return c.Close()
}

func (c *BaseClient) OnStart() func(dispatch.StartCallback) dispatch.StartCallback {
	return c.router.OnStart()
}

func (c *BaseClient) OnMessage(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnMessage(filters...)
}

func (c *BaseClient) OnMessageEdit(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnMessageEdit(filters...)
}

func (c *BaseClient) OnMessageDelete(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnMessageDelete(filters...)
}

func (c *BaseClient) OnMessageRead(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnMessageRead(filters...)
}

func (c *BaseClient) OnTyping(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnTyping(filters...)
}

func (c *BaseClient) OnPresence(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnPresence(filters...)
}

func (c *BaseClient) OnReactionUpdate(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnReactionUpdate(filters...)
}

func (c *BaseClient) OnChatUpdate(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnChatUpdate(filters...)
}

func (c *BaseClient) OnRaw(filters ...dispatch.FilterCallback) func(dispatch.HandlerCallback) dispatch.HandlerCallback {
	return c.router.OnRaw(filters...)
}

func (c *BaseClient) OnError(scope ErrorScope) func(dispatch.ErrorCallback) dispatch.ErrorCallback {
	return c.router.OnError(dispatch.ErrorScope(scope))
}

func (c *BaseClient) OnDisconnect() func(dispatch.DisconnectCallback) dispatch.DisconnectCallback {
	return c.router.OnDisconnect()
}

func (c *BaseClient) IncludeRouter(router *dispatch.Router) {
	c.router.IncludeRouter(router)
}

func (c *BaseClient) Relogin(dropConfigToken bool, start bool) error {
	if c.app == nil || c.app.session == nil {
		return fmt.Errorf("cannot relogin before session is loaded")
	}

	c.app.store.DeleteSession(c.app.session.Token)
	c.Close()

	if dropConfigToken {
		c.ExtraConfig.Token = ""
		c.config.Token = ""
	}

	c.ResetRuntime()
	if start {
		return c.Start()
	}
	return nil
}

func generateUUID() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(time.Now().UnixNano()>>uint(i*4)) ^ byte(i*13)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
