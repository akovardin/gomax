package gomax

import (
	"github.com/akovardin/gomax/auth"
	"github.com/akovardin/gomax/connection"
	"github.com/akovardin/gomax/connection/readers"
	ws "github.com/akovardin/gomax/protocol/ws"
	"github.com/akovardin/gomax/transport"
)

type WebClient struct {
	BaseClient

	qrProvider       auth.QrHandler
	passwordProvider auth.PasswordProvider
}

func NewWebClient(
	sessionName string,
	workDir string,
	extraConfig *ExtraConfig,
	qrProvider auth.QrHandler,
	passwordProvider auth.PasswordProvider,
) *WebClient {
	if extraConfig == nil {
		extraConfig = DefaultExtraConfig()
	}
	if qrProvider == nil {
		qrProvider = &auth.ConsoleQrHandler{}
	}
	if passwordProvider == nil {
		passwordProvider = &auth.ConsolePasswordProvider{}
	}

	client := &WebClient{
		qrProvider:       qrProvider,
		passwordProvider: passwordProvider,
	}
	client.SessionName = sessionName
	client.WorkDir = workDir
	client.ExtraConfig = extraConfig

	userAgent := extraConfig.GenerateWebUserAgent()
	client.config = client.BuildConfig("", userAgent)

	client.buildConn = client.buildConnection
	flow := auth.NewQrAuthFlow(qrProvider, passwordProvider)
	client.InitRuntime(flow)

	return client
}

func (c *WebClient) buildConnection() *connection.ConnectionManager {
	url := c.ExtraConfig.GetURL()
	trans := transport.NewWebsocketTransport(url, c.ExtraConfig.Proxy)
	proto := ws.NewWsProtocol()
	reader := readers.NewWsReader(trans)

	return connection.NewConnectionManager(reader, trans, proto)
}

func (c *WebClient) AuthFlow() *auth.QrAuthFlow {
	return auth.NewQrAuthFlow(c.qrProvider, c.passwordProvider)
}
