package auth

import (
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/types"
)

type Service struct {
	app core.AppInterface
}

func NewService(app core.AppInterface) *Service {
	return &Service{app: app}
}

func (s *Service) resolveSyncState() types.SyncState {
	config := s.app.Config()
	saved := types.DefaultSyncState()
	if config.Sync != nil {
		saved = config.Sync.Resolve(saved)
	}
	return saved
}

func (s *Service) RequestCode(phone string) (*types.StartAuthResponse, error) {
	core.LogDebug("requestCode phone=%s", phone)
	payload := map[string]interface{}{
		"phone":    phone,
		"type":     "START_AUTH",
		"language": "ru",
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuthRequest), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.StartAuthResponse](frame.Payload)
}

func (s *Service) SendCode(token string, code string) (*types.CheckCodeResponse, error) {
	tokenPrefix := token
	if len(token) > 8 {
		tokenPrefix = token[:8]
	}
	core.LogDebug("sendCode token=%s...", tokenPrefix)
	payload := map[string]interface{}{
		"token":          token,
		"verifyCode":     code,
		"authTokenType":  "CHECK_CODE",
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuth), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.CheckCodeResponse](frame.Payload)
}

func (s *Service) CheckPassword(trackID string, password string) (*types.CheckPasswordResponse, error) {
	payload := map[string]interface{}{
		"trackId":  trackID,
		"password": password,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuthLoginCheckPassword), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.CheckPasswordResponse](frame.Payload)
}

func (s *Service) Login(userAgent *types.MobileUserAgentPayload) (*types.LoginResponse, error) {
	core.LogDebug("login deviceType=%s", userAgent.DeviceType)
	if userAgent != nil && userAgent.DeviceType == string(types.DeviceTypeWeb) {
		return s.WebLogin()
	}
	return s.MobileLogin()
}

func (s *Service) MobileLogin() (*types.LoginResponse, error) {
	syncState := s.resolveSyncState()
	config := s.app.Config()

	payload := map[string]interface{}{
		"userAgent":       config.Device.UserAgent,
		"token":           s.app.Token(),
		"chatsSync":       syncState.ChatsSync,
		"contactsSync":    syncState.ContactsSync,
		"draftsSync":      syncState.DraftsSync,
		"presenceSync":    syncState.PresenceSync,
		"interactive":     true,
		"configHash":      syncState.ConfigHash,
		"deviceId":        config.Device.DeviceID,
		"mtInstanceId":    config.Device.MtInstanceID,
		"clientSessionId": config.Device.ClientSessionID,
	}

	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeLogin), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.LoginResponse](frame.Payload)
}

func (s *Service) WebLogin() (*types.LoginResponse, error) {
	syncState := s.resolveSyncState()

	payload := map[string]interface{}{
		"token":        s.app.Token(),
		"chatsCount":   40,
		"interactive":  true,
		"chatsSync":    syncState.ChatsSync,
		"contactsSync": syncState.ContactsSync,
		"presenceSync": syncState.PresenceSync,
		"draftsSync":   syncState.DraftsSync,
	}

	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeLogin), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.LoginResponse](frame.Payload)
}

func (s *Service) RequestQr() (*types.RequestQrResponse, error) {
	core.LogDebug("requestQr")
	payload := map[string]interface{}{}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeGetQR), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.RequestQrResponse](frame.Payload)
}

func (s *Service) CheckQr(trackID string) (*types.CheckQrResponse, error) {
	core.LogDebug("checkQr trackId=%s", trackID)
	payload := map[string]interface{}{
		"trackId": trackID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeGetQRStatus), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.CheckQrResponse](frame.Payload)
}

func (s *Service) ConfirmQr(trackID string) (*types.CheckCodeResponse, error) {
	core.LogDebug("confirmQr trackId=%s", trackID)
	payload := map[string]interface{}{
		"trackId": trackID,
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeLoginByQR), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.CheckCodeResponse](frame.Payload)
}

func (s *Service) ConfirmRegistration(firstName string, lastName string, token string) (*types.ConfirmRegistrationResponse, error) {
	payload := map[string]interface{}{
		"firstName": firstName,
		"lastName":  lastName,
		"token":     token,
		"tokenType": "REGISTER",
	}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuthConfirm), payload)
	if err != nil {
		return nil, err
	}
	return core.RequirePayloadModel[types.ConfirmRegistrationResponse](frame.Payload)
}

func (s *Service) Set2FA(password string, email string, hint string) (bool, error) {
	payload := map[string]interface{}{
		"password": password,
		"email":    email,
		"hint":     hint,
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuthSet2FA), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Remove2FA(password string) (bool, error) {
	payload := map[string]interface{}{
		"password": password,
		"action":   "remove",
	}
	_, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuthSet2FA), payload)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Check2FA() (bool, error) {
	payload := map[string]interface{}{}
	frame, err := core.InvokeAPI(s.app, int(protocol.OpcodeAuth2FADetails), payload)
	if err != nil {
		return false, err
	}
	if enabled, ok := frame.Payload["enabled"].(bool); ok {
		return enabled, nil
	}
	return false, nil
}
