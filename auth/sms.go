package auth

import (
	"fmt"

	apiauth "github.com/akovardin/gomax/api/auth"
	"github.com/akovardin/gomax/api/core"
)

type SmsAuthFlow struct {
	CodeProvider     SmsCodeProvider
	PasswordProvider PasswordProvider
}

func NewSmsAuthFlow(codeProvider SmsCodeProvider, passwordProvider PasswordProvider) *SmsAuthFlow {
	return &SmsAuthFlow{
		CodeProvider:     codeProvider,
		PasswordProvider: passwordProvider,
	}
}

func (f *SmsAuthFlow) Authenticate(app core.AppInterface) (*AuthResult, error) {
	config := app.Config()
	phone := config.Phone
	if phone == "" {
		return nil, fmt.Errorf("phone is required for SMS authentication")
	}

	core.LogDebug("sms auth: requesting code for %s", phone)

	authSvc := apiauth.NewService(app)

	start, err := authSvc.RequestCode(phone)
	if err != nil {
		core.LogWarn("sms code error: %v", err)
		return nil, err
	}

	core.LogDebug("sms auth: code received, sending...")

	code, err := f.CodeProvider.GetCode(phone)
	if err != nil {
		core.LogWarn("sms get code error: %v", err)
		return nil, err
	}

	result, err := authSvc.SendCode(start.Token, code)
	if err != nil {
		core.LogWarn("sms send code error: %v", err)
		return nil, err
	}

	var token string
	if lt := result.TokenAttrs.LoginToken(); lt != nil {
		core.LogDebug("sms auth: got login token")
		token = *lt
	} else if result.PasswordChallenge != nil {
		core.LogDebug("sms auth: password challenge")
		token, err = f.authenticateWithPassword(app, result.PasswordChallenge.TrackID, result.PasswordChallenge.Hint)
		if err != nil {
			return nil, err
		}
	} else if rt := result.TokenAttrs.RegisterToken(); rt != nil {
		core.LogDebug("sms auth: registration required")
		registrationConfig := config.RegistrationConfig
		if registrationConfig == nil {
			return nil, fmt.Errorf("RegistrationConfig is required to register a new account")
		}
		response, err := authSvc.ConfirmRegistration(registrationConfig.FirstName, registrationConfig.LastName, *rt)
		if err != nil {
			return nil, err
		}
		token = response.Token
	} else {
		return nil, fmt.Errorf("authentication failed: no login token, password challenge, or registration token")
	}

	return &AuthResult{Token: token}, nil
}

func (f *SmsAuthFlow) authenticateWithPassword(app core.AppInterface, trackID string, hint *string) (string, error) {
	authSvc := apiauth.NewService(app)
	hintStr := ""
	if hint != nil {
		hintStr = *hint
	}

	for {
		password, err := f.PasswordProvider.GetPassword(hintStr)
		if err != nil {
			return "", err
		}
		if password == "" {
			continue
		}

		response, err := authSvc.CheckPassword(trackID, password)
		if err != nil {
			continue
		}

		if response.Error != nil && *response.Error != "" {
			continue
		}

		if lt := response.LoginToken(); lt != nil {
			return *lt, nil
		}
	}
}
