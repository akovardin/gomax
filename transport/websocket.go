package transport

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketTransport struct {
	url   string
	proxy string
	conn  *websocket.Conn
}

func NewWebsocketTransport(url string, proxy string) *WebsocketTransport {
	return &WebsocketTransport{url: url, proxy: proxy}
}

func (t *WebsocketTransport) Connect() error {
	header := http.Header{}
	header.Set("Origin", "https://web.max.ru")

	conn, _, err := websocket.DefaultDialer.Dial(t.url, header)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *WebsocketTransport) Close() error {
	if t.conn == nil {
		return nil
	}
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	t.conn.WriteMessage(websocket.CloseMessage, msg)
	return t.conn.Close()
}

func (t *WebsocketTransport) Send(data []byte) error {
	if t.conn == nil {
		return websocket.ErrCloseSent
	}
	return t.conn.WriteMessage(websocket.TextMessage, data)
}

func (t *WebsocketTransport) Recv() ([]byte, error) {
	if t.conn == nil {
		return nil, websocket.ErrCloseSent
	}
	_, msg, err := t.conn.ReadMessage()
	return msg, err
}

func (t *WebsocketTransport) Connected() bool {
	return t.conn != nil
}
