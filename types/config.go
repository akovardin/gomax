package types

type DeviceType string

const (
	DeviceTypeWeb     DeviceType = "WEB"
	DeviceTypeAndroid DeviceType = "ANDROID"
	DeviceTypeIOS     DeviceType = "IOS"
	DeviceTypeDesktop DeviceType = "DESKTOP"
)

type MobileUserAgentPayload struct {
	DeviceType      string `json:"deviceType" msgpack:"deviceType"`
	AppVersion      string `json:"appVersion" msgpack:"appVersion"`
	OsVersion       string `json:"osVersion" msgpack:"osVersion"`
	Timezone        string `json:"timezone" msgpack:"timezone"`
	Screen          string `json:"screen" msgpack:"screen"`
	PushDeviceType  string `json:"pushDeviceType" msgpack:"pushDeviceType"`
	Arch            string `json:"arch" msgpack:"arch"`
	Locale          string `json:"locale" msgpack:"locale"`
	BuildNumber     int    `json:"buildNumber" msgpack:"buildNumber"`
	DeviceName      string `json:"deviceName" msgpack:"deviceName"`
	DeviceLocale    string `json:"deviceLocale" msgpack:"deviceLocale"`
	HeaderUserAgent string `json:"headerUserAgent" msgpack:"headerUserAgent"`
}

func (p *MobileUserAgentPayload) ToWebPayload() map[string]interface{} {
	return map[string]interface{}{
		"deviceType":      p.DeviceType,
		"locale":          p.Locale,
		"deviceLocale":    p.DeviceLocale,
		"osVersion":       p.OsVersion,
		"deviceName":      p.DeviceName,
		"headerUserAgent": p.HeaderUserAgent,
		"appVersion":      p.AppVersion,
		"screen":          p.Screen,
		"timezone":        p.Timezone,
	}
}

type DeviceConfig struct {
	MtInstanceID    string                  `json:"mtInstanceId"`
	UserAgent       *MobileUserAgentPayload `json:"userAgent"`
	DeviceID        string                  `json:"deviceId"`
	ClientSessionID int                     `json:"clientSessionId"`
}

type RegistrationConfig struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ClientConfig struct {
	Phone              string              `json:"phone"`
	WorkDir            string              `json:"workDir"`
	SessionName        string              `json:"sessionName"`
	Device             *DeviceConfig       `json:"device"`
	Token              string              `json:"token"`
	Proxy              string              `json:"proxy"`
	RegistrationConfig *RegistrationConfig `json:"registrationConfig"`
	Host               string              `json:"host"`
	Port               int                 `json:"port"`
	UseSSL             bool                `json:"useSsl"`
	ProtocolVersion    int                 `json:"protocolVersion"`
	RequestTimeout     float64             `json:"requestTimeout"`
	LogLevel           string              `json:"logLevel"`
	Telemetry          bool                `json:"telemetry"`
	Sync               *SyncOverrides      `json:"sync"`
}
