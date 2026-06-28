package auth

import "github.com/akovardin/gomax/api/core"

type AuthResult struct {
	Token string
}

type AuthFlow interface {
	Authenticate(app core.AppInterface) (*AuthResult, error)
}
