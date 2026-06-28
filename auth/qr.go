package auth

import (
	"fmt"
	"time"

	apiauth "github.com/akovardin/gomax/api/auth"
	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/types"
)

type QrAuthFlow struct {
	QRHandler        QrHandler
	PasswordProvider PasswordProvider
}

func NewQrAuthFlow(qrHandler QrHandler, passwordProvider PasswordProvider) *QrAuthFlow {
	return &QrAuthFlow{
		QRHandler:        qrHandler,
		PasswordProvider: passwordProvider,
	}
}

func (f *QrAuthFlow) Authenticate(app core.AppInterface) (*AuthResult, error) {
	authSvc := apiauth.NewService(app)

	qrInfo, err := authSvc.RequestQr()
	if err != nil {
		return nil, err
	}

	interval := float64(qrInfo.PollingInterval)
	expiresAt := float64(qrInfo.ExpiresAt)
	core.LogDebug("pollQR trackId=%s interval=%.1fs expiresIn=%ds", qrInfo.TrackID, interval/1000.0, int64(expiresAt/1000))

	if err := f.QRHandler.ShowQR(qrInfo.QRLink); err != nil {
		return nil, err
	}

	confirmed, err := f.pollQR(app, qrInfo)
	if err != nil {
		return nil, err
	}
	if !confirmed {
		return nil, fmt.Errorf("QR authentication expired")
	}

	result, err := authSvc.ConfirmQr(qrInfo.TrackID)
	if err != nil {
		return nil, err
	}

	var token string
	if lt := result.TokenAttrs.LoginToken(); lt != nil {
		token = *lt
	} else if result.PasswordChallenge != nil {
		token, err = f.authenticateWithPassword(app, result.PasswordChallenge.TrackID, result.PasswordChallenge.Hint)
		if err != nil {
			return nil, err
		}
	}

	return &AuthResult{Token: token}, nil
}

func (f *QrAuthFlow) pollQR(app core.AppInterface, qrInfo *types.RequestQrResponse) (bool, error) {
	authSvc := apiauth.NewService(app)
	interval := float64(qrInfo.PollingInterval) / 1000.0
	expiresAt := float64(qrInfo.ExpiresAt) / 1000.0

	for float64(time.Now().Unix()) < expiresAt {
		response, err := authSvc.CheckQr(qrInfo.TrackID)
		if err != nil {
			return false, err
		}

		if response.Status.LoginAvailable {
			core.LogDebug("qr confirmed")
			return true, nil
		}

		time.Sleep(time.Duration(interval * float64(time.Second)))
	}

	return false, nil
}

func (f *QrAuthFlow) authenticateWithPassword(app core.AppInterface, trackID string, hint *string) (string, error) {
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
