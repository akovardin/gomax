package types

type Session struct {
	ID           FlexInt                `json:"id"`
	DeviceID     string                 `json:"deviceId"`
	Current      bool                   `json:"current"`
	UserAgent    string                 `json:"userAgent"`
	AppVersion   string                 `json:"appVersion"`
	DeviceName   string                 `json:"deviceName"`
	DeviceType   string                 `json:"deviceType"`
	Platform     string                 `json:"platform"`
	IP           string                 `json:"ip"`
	Location     string                 `json:"location"`
	Created      FlexInt                `json:"created"`
	Updated      FlexInt                `json:"updated"`
	LastActivity FlexInt                `json:"lastActivity"`
	Options      map[string]interface{} `json:"options"`
}
