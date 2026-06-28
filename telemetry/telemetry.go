package telemetry

import (
	"math/rand"
	"time"

	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/protocol"
)

// Screen represents a virtual app screen
type Screen string

const (
	ScreenBackground Screen = "BACKGROUND"
	ScreenChats      Screen = "CHATS"
	ScreenChat       Screen = "CHAT"
	ScreenSettings   Screen = "SETTINGS"
	ScreenProfile    Screen = "PROFILE"
	ScreenContacts   Screen = "CONTACTS"
	ScreenCalls      Screen = "CALLS"
)

// TelemetryService sends fake app usage data.
type TelemetryService struct {
	app           core.AppInterface
	stopCh        chan struct{}
	currentScreen Screen
}

func NewTelemetryService(app core.AppInterface) *TelemetryService {
	return &TelemetryService{
		app:           app,
		stopCh:        make(chan struct{}),
		currentScreen: ScreenBackground,
	}
}

// Start begins sending telemetry events in the background.
func (t *TelemetryService) Start() {
	go t.loop()
}

// Stop stops sending telemetry events.
func (t *TelemetryService) Stop() {
	close(t.stopCh)
}

func (t *TelemetryService) loop() {
	// Wait a bit before starting
	time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)

	for {
		select {
		case <-t.stopCh:
			return
		default:
		}

		t.navigate()
		t.sendPerf()

		// Random delay between events
		delay := time.Duration(5+rand.Intn(25)) * time.Second
		select {
		case <-t.stopCh:
			return
		case <-time.After(delay):
		}
	}
}

func (t *TelemetryService) navigate() {
	// Simple navigation simulation: cycle through common screens
	nextScreen := t.pickNextScreen()

	t.sendNavEvent("GO", string(t.currentScreen), string(nextScreen))
	t.currentScreen = nextScreen
}

func (t *TelemetryService) pickNextScreen() Screen {
	// Weighted random selection
	screens := []Screen{ScreenChats, ScreenChat, ScreenContacts, ScreenSettings, ScreenProfile}
	return screens[rand.Intn(len(screens))]
}

func (t *TelemetryService) sendNavEvent(action string, from string, to string) {
	payload := map[string]interface{}{
		"type":   "NAV",
		"action": action,
		"from":   from,
		"to":     to,
		"time":   time.Now().UnixMilli(),
	}

	t.app.Invoke(
		int(protocol.OpcodeLog),
		payload,
		int(protocol.CommandRequest),
		5.0,
		false,
	)
}

func (t *TelemetryService) sendPerf() {
	events := []string{"open_chat_to_render", "open_chats_to_render", "login", "sync_contacts"}
	event := events[rand.Intn(len(events))]

	payload := map[string]interface{}{
		"type":     "PERF",
		"event":    event,
		"duration": rand.Intn(2000),
		"time":     time.Now().UnixMilli(),
	}

	t.app.Invoke(
		int(protocol.OpcodeLog),
		payload,
		int(protocol.CommandRequest),
		5.0,
		false,
	)
}
