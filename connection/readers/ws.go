package readers

import "github.com/akovardin/gomax/transport"

type WsReader struct {
	transport transport.Transport
}

func NewWsReader(transport transport.Transport) *WsReader {
	return &WsReader{transport: transport}
}

func (r *WsReader) Read() ([]byte, error) {
	return r.transport.Recv()
}
