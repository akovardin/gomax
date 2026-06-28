package gomax

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand/v2"

	"github.com/akovardin/gomax/types"
)

const (
	DefaultHost          = "api.oneme.ru"
	DefaultPort          = 443
	DefaultWSURL         = "wss://ws-api.oneme.ru/websocket"
	DefaultProtocolVersion = 1
	DefaultRequestTimeout = 30.0
	DefaultLogLevel       = "info"
)

type appVersionEntry struct {
	AppVersion  string
	BuildNumber int
	OsVersion   string
	Arch        string
}

type androidDeviceEntry struct {
	DeviceName     string
	DeviceType     string
	Screen         string
	PushDeviceType string
}

type localeTimezoneEntry struct {
	Locale   string
	Timezone string
}

var appVersions = []appVersionEntry{
	{AppVersion: "26.14.1", BuildNumber: 6686, OsVersion: "14", Arch: "arm64-v8a"},
	{AppVersion: "26.14.0", BuildNumber: 6685, OsVersion: "14", Arch: "arm64-v8a"},
	{AppVersion: "26.13.0", BuildNumber: 6683, OsVersion: "14", Arch: "arm64-v8a"},
	{AppVersion: "26.12.2", BuildNumber: 6681, OsVersion: "14", Arch: "arm64-v8a"},
	{AppVersion: "26.12.1", BuildNumber: 6679, OsVersion: "13", Arch: "arm64-v8a"},
	{AppVersion: "26.12.0", BuildNumber: 6678, OsVersion: "13", Arch: "arm64-v8a"},
	{AppVersion: "26.11.3", BuildNumber: 6680, OsVersion: "13", Arch: "arm64-v8a"},
	{AppVersion: "26.11.2", BuildNumber: 6669, OsVersion: "13", Arch: "arm64-v8a"},
	{AppVersion: "26.11.1", BuildNumber: 6665, OsVersion: "13", Arch: "arm64-v8a"},
	{AppVersion: "26.11.0", BuildNumber: 6661, OsVersion: "13", Arch: "arm64-v8a"},
}

var androidDevices = []androidDeviceEntry{
	{DeviceName: "Samsung SM-A525F", DeviceType: "ANDROID", Screen: "405dpi 405dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Samsung SM-A536B", DeviceType: "ANDROID", Screen: "405dpi 405dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Samsung SM-A546E", DeviceType: "ANDROID", Screen: "405dpi 405dpi 1080x2340", PushDeviceType: "GCM"},
	{DeviceName: "Samsung SM-G991B", DeviceType: "ANDROID", Screen: "421dpi 421dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Samsung SM-G998B", DeviceType: "ANDROID", Screen: "515dpi 515dpi 1440x3200", PushDeviceType: "GCM"},
	{DeviceName: "Xiaomi 2109119DG", DeviceType: "ANDROID", Screen: "395dpi 395dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Xiaomi 2201117TG", DeviceType: "ANDROID", Screen: "395dpi 395dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Pixel 6", DeviceType: "ANDROID", Screen: "411dpi 411dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Pixel 7", DeviceType: "ANDROID", Screen: "416dpi 416dpi 1080x2400", PushDeviceType: "GCM"},
	{DeviceName: "Pixel 8", DeviceType: "ANDROID", Screen: "428dpi 428dpi 1080x2400", PushDeviceType: "GCM"},
}

var localeTimezones = []localeTimezoneEntry{
	{Locale: "ru_RU", Timezone: "Europe/Moscow"},
	{Locale: "en_US", Timezone: "America/New_York"},
	{Locale: "uk_UA", Timezone: "Europe/Kiev"},
	{Locale: "en_GB", Timezone: "Europe/London"},
	{Locale: "de_DE", Timezone: "Europe/Berlin"},
	{Locale: "fr_FR", Timezone: "Europe/Paris"},
	{Locale: "es_ES", Timezone: "Europe/Madrid"},
	{Locale: "it_IT", Timezone: "Europe/Rome"},
	{Locale: "pt_BR", Timezone: "America/Sao_Paulo"},
	{Locale: "tr_TR", Timezone: "Europe/Istanbul"},
	{Locale: "kk_KZ", Timezone: "Asia/Almaty"},
	{Locale: "be_BY", Timezone: "Europe/Minsk"},
	{Locale: "uz_UZ", Timezone: "Asia/Tashkent"},
	{Locale: "ja_JP", Timezone: "Asia/Tokyo"},
	{Locale: "ko_KR", Timezone: "Asia/Seoul"},
}

var rng *rand.Rand

func init() {
	var seed [8]byte
	if _, err := cryptorand.Read(seed[:]); err != nil {
		rng = rand.New(rand.NewPCG(uint64(42), uint64(17)))
		return
	}
	rng = rand.New(rand.NewPCG(binary.LittleEndian.Uint64(seed[:]), binary.LittleEndian.Uint64(seed[:])))
}

type ExtraConfig struct {
	Token              string                       `json:"token"`
	RegistrationConfig *types.RegistrationConfig    `json:"registrationConfig"`
	Host               string                       `json:"host"`
	Port               int                          `json:"port"`
	URL                string                       `json:"url"`
	UseSSL             bool                         `json:"useSsl"`
	Proxy              string                       `json:"proxy"`
	Reconnect          bool                         `json:"reconnect"`
	ReconnectDelay     float64                      `json:"reconnectDelay"`
	DeviceID           string                       `json:"deviceId"`
	DeviceType         types.DeviceType             `json:"deviceType"`
	UserAgent          *types.MobileUserAgentPayload `json:"userAgent"`
	MtInstanceID       string                       `json:"mtInstanceId"`
	RequestTimeout     float64                      `json:"requestTimeout"`
	LogLevel           string                       `json:"logLevel"`
	Telemetry          bool                         `json:"telemetry"`
	Sync               *types.SyncOverrides         `json:"sync"`
}

func (c *ExtraConfig) GenerateUserAgent() *types.MobileUserAgentPayload {
	if c.UserAgent != nil {
		return c.UserAgent
	}

	av := appVersions[rng.IntN(len(appVersions))]
	dev := androidDevices[rng.IntN(len(androidDevices))]
	lt := localeTimezones[rng.IntN(len(localeTimezones))]

	deviceLang := lt.Locale[:2]

	headerUA := fmt.Sprintf(
		"Max-Android/%s (Linux; Android %s; %s) MaxApp/%s (Android; %s; %s; %s; %s; %s)",
		av.AppVersion,
		av.OsVersion,
		dev.DeviceName,
		av.AppVersion,
		av.OsVersion,
		dev.Screen,
		lt.Locale,
		av.Arch,
		lt.Timezone,
	)

	return &types.MobileUserAgentPayload{
		DeviceType:      dev.DeviceType,
		AppVersion:      av.AppVersion,
		OsVersion:       av.OsVersion,
		Timezone:        lt.Timezone,
		Screen:          dev.Screen,
		PushDeviceType:  dev.PushDeviceType,
		Arch:            av.Arch,
		Locale:          lt.Locale,
		BuildNumber:     av.BuildNumber,
		DeviceName:      dev.DeviceName,
		DeviceLocale:    deviceLang,
		HeaderUserAgent: headerUA,
	}
}

const (
	webAppVersion = "26.5.5"
	webScreen     = "1080x1920 1.0x"
	defaultWebHeaderUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"
)

func (c *ExtraConfig) GenerateWebUserAgent() *types.MobileUserAgentPayload {
	if c.UserAgent != nil {
		return c.UserAgent
	}

	lt := localeTimezones[rng.IntN(len(localeTimezones))]

	return &types.MobileUserAgentPayload{
		DeviceType:      string(types.DeviceTypeWeb),
		AppVersion:      webAppVersion,
		OsVersion:       "Linux",
		Timezone:        lt.Timezone,
		Screen:          webScreen,
		Locale:          lt.Locale,
		DeviceName:      "Chrome",
		DeviceLocale:    lt.Locale[:2],
		HeaderUserAgent: defaultWebHeaderUserAgent,
	}
}

func (c *ExtraConfig) GetURL() string {
	if c.URL != "" {
		return c.URL
	}
	scheme := "ws"
	if c.UseSSL {
		scheme = "wss"
	}
	host := c.Host
	if host == "" {
		host = DefaultHost
	}
	port := c.Port
	if port == 0 {
		port = DefaultPort
	}
	return fmt.Sprintf("%s://%s:%d/websocket", scheme, host, port)
}

func DefaultExtraConfig() *ExtraConfig {
	return &ExtraConfig{
		Host:           DefaultHost,
		Port:           DefaultPort,
		URL:            DefaultWSURL,
		UseSSL:         true,
		Reconnect:      true,
		ReconnectDelay: 5.0,
		RequestTimeout: DefaultRequestTimeout,
		LogLevel:       DefaultLogLevel,
		Telemetry:      true,
	}
}

func DefaultClientConfig() *types.ClientConfig {
	return &types.ClientConfig{
		Host:            DefaultHost,
		Port:            DefaultPort,
		UseSSL:          true,
		ProtocolVersion: DefaultProtocolVersion,
		RequestTimeout:  DefaultRequestTimeout,
		LogLevel:        DefaultLogLevel,
	}
}
