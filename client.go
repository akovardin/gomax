package gomax

import (
	"github.com/akovardin/gomax/auth"
	"github.com/akovardin/gomax/connection"
	"github.com/akovardin/gomax/connection/readers"
	tcp "github.com/akovardin/gomax/protocol/tcp"
	"github.com/akovardin/gomax/transport"
	"github.com/akovardin/gomax/types"
)

type Client struct {
	BaseClient

	smsCodeProvider  auth.SmsCodeProvider
	passwordProvider auth.PasswordProvider
}

func NewClient(
	phone string,
	sessionName string,
	workDir string,
	extraConfig *ExtraConfig,
	smsCodeProvider auth.SmsCodeProvider,
	passwordProvider auth.PasswordProvider,
) *Client {
	if extraConfig == nil {
		extraConfig = DefaultExtraConfig()
	}
	if smsCodeProvider == nil {
		smsCodeProvider = &auth.ConsoleSmsCodeProvider{}
	}
	if passwordProvider == nil {
		passwordProvider = &auth.ConsolePasswordProvider{}
	}

	client := &Client{
		smsCodeProvider:  smsCodeProvider,
		passwordProvider: passwordProvider,
	}
	client.Phone = phone
	client.SessionName = sessionName
	client.WorkDir = workDir
	client.ExtraConfig = extraConfig

	userAgent := extraConfig.GenerateUserAgent()
	client.config = client.BuildConfig(phone, userAgent)

	client.buildConn = client.buildConnection
	flow := auth.NewSmsAuthFlow(smsCodeProvider, passwordProvider)
	client.InitRuntime(flow)

	return client
}

func (c *Client) buildConnection() *connection.ConnectionManager {
	trans := transport.NewTcpTransport(
		c.ExtraConfig.Host,
		c.ExtraConfig.Port,
		c.ExtraConfig.Proxy,
		c.ExtraConfig.UseSSL,
	)
	proto := tcp.NewTcpProtocol()
	reader := readers.NewTcpReader(trans, proto.GetFramer())

	return connection.NewConnectionManager(reader, trans, proto)
}

func (c *Client) AuthFlow() *auth.SmsAuthFlow {
	return auth.NewSmsAuthFlow(c.smsCodeProvider, c.passwordProvider)
}

func (c *Client) WithExtraConfig(extraConfig *ExtraConfig) *Client {
	c.ExtraConfig = extraConfig
	configPhone := c.config.Phone
	if c.Phone != "" {
		configPhone = c.Phone
	}
	userAgent := extraConfig.GenerateUserAgent()
	c.config = c.BuildConfig(configPhone, userAgent)
	return c
}

func (c *Client) WithUserAgent(userAgent *types.MobileUserAgentPayload) *Client {
	c.config.Device.UserAgent = userAgent
	return c
}

func (c *Client) BuildConfig(phone string, userAgent *types.MobileUserAgentPayload) *types.ClientConfig {
	return c.BaseClient.BuildConfig(phone, userAgent)
}
