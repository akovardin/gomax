package session

import (
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/logging"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
)

type Service struct {
	app core.AppInterface
}

func NewService(app core.AppInterface) *Service {
	return &Service{app: app}
}

func (s *Service) Handshake(mtInstanceID string, userAgent *types.MobileUserAgentPayload, deviceID string) error {
	if userAgent != nil && userAgent.DeviceType == string(types.DeviceTypeWeb) {
		return s.webHandshake(userAgent, deviceID)
	}
	return s.mobileHandshake(mtInstanceID, userAgent, deviceID)
}

func (s *Service) mobileHandshake(mtInstanceID string, userAgent *types.MobileUserAgentPayload, deviceID string) error {
	logging.LogDebug("handshake deviceId=%s deviceType=%s", deviceID, userAgent.DeviceType)
	payload := map[string]interface{}{
		"deviceId":     deviceID,
		"mtInstanceId": mtInstanceID,
	}
	if userAgent != nil {
		payload["userAgent"] = userAgent
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeSessionInit), payload)
	return err
}

func (s *Service) webHandshake(userAgent *types.MobileUserAgentPayload, deviceID string) error {
	logging.LogDebug("handshake deviceId=%s deviceType=%s", deviceID, userAgent.DeviceType)
	payload := map[string]interface{}{
		"deviceId": deviceID,
	}
	if userAgent != nil {
		payload["userAgent"] = userAgent.ToWebPayload()
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeSessionInit), payload)
	return err
}
