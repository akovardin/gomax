package types

type StartAuthResponse struct {
	Token              string  `json:"token"`
	CodeLength         FlexInt `json:"codeLength"`
	RequestMaxDuration FlexInt `json:"requestMaxDuration"`
	RequestCountLeft   FlexInt `json:"requestCountLeft"`
	AltActionDuration  FlexInt `json:"altActionDuration"`
}

type Token struct {
	Token string `json:"token"`
}

type TokenAttrs struct {
	Login    *Token `json:"LOGIN,omitempty"`
	Register *Token `json:"REGISTER,omitempty"`
}

func (t *TokenAttrs) LoginToken() *string {
	if t != nil && t.Login != nil {
		return &t.Login.Token
	}
	return nil
}

func (t *TokenAttrs) RegisterToken() *string {
	if t != nil && t.Register != nil {
		return &t.Register.Token
	}
	return nil
}

type PasswordChallenge struct {
	TrackID string  `json:"trackId"`
	Hint    *string `json:"hint,omitempty"`
}

type CheckCodeResponse struct {
	TokenAttrs        *TokenAttrs        `json:"tokenAttrs,omitempty"`
	PasswordChallenge *PasswordChallenge `json:"passwordChallenge,omitempty"`
}

type CheckPasswordResponse struct {
	TokenAttrs *TokenAttrs `json:"tokenAttrs,omitempty"`
	Error      *string     `json:"error,omitempty"`
}

func (c *CheckPasswordResponse) LoginToken() *string {
	return c.TokenAttrs.LoginToken()
}

type QrStatus struct {
	ExpiresAt      FlexInt `json:"expiresAt"`
	LoginAvailable bool    `json:"loginAvailable"`
}

type RequestQrResponse struct {
	ExpiresAt       FlexInt `json:"expiresAt"`
	PollingInterval FlexInt `json:"pollingInterval"`
	QRLink          string  `json:"qrLink"`
	TrackID         string  `json:"trackId"`
	TTL             FlexInt `json:"ttl"`
}

type CheckQrResponse struct {
	Status QrStatus `json:"status"`
}

type ConfirmRegistrationResponse struct {
	UserToken string   `json:"userToken"`
	Profile   *Profile `json:"profile"`
	TokenType FlexInt  `json:"tokenType"`
	Token     string   `json:"token"`
}
