package ws

import (
	"encoding/json"

	"github.com/akovardin/gomax/protocol"
)

type WsProtocol struct{}

func NewWsProtocol() *WsProtocol {
	return &WsProtocol{}
}

func (p *WsProtocol) Version() int {
	return 11
}

func (p *WsProtocol) Encode(frame *protocol.OutboundFrame) ([]byte, error) {
	return json.Marshal(frame)
}

func (p *WsProtocol) Decode(raw []byte) (*protocol.InboundFrame, error) {
	frame := &protocol.InboundFrame{}
	if err := json.Unmarshal(raw, frame); err != nil {
		return nil, err
	}
	return frame, nil
}
