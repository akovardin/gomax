package readers

import (
	"fmt"

	"github.com/akovardin/gomax/protocol/tcp"
	"github.com/akovardin/gomax/transport"
)

type TcpReader struct {
	transport transport.Transport
	framer    *tcp.TcpPacketFramer
	buf       []byte
}

func NewTcpReader(transport transport.Transport, framer *tcp.TcpPacketFramer) *TcpReader {
	return &TcpReader{transport: transport, framer: framer}
}

const tcpHeaderSize = 10

func (r *TcpReader) Read() ([]byte, error) {
	header, err := r.readExact(tcpHeaderSize)
	if err != nil {
		return nil, err
	}

	payloadLen, ok := r.framer.UnpackHeader(header)
	if !ok {
		return nil, fmt.Errorf("tcp: incomplete header")
	}

	payload, err := r.readExact(payloadLen)
	if err != nil {
		return nil, err
	}

	result := make([]byte, tcpHeaderSize+len(payload))
	copy(result, header)
	copy(result[tcpHeaderSize:], payload)
	return result, nil
}

func (r *TcpReader) readExact(n int) ([]byte, error) {
	result := make([]byte, n)
	offset := 0
	for offset < n {
		if len(r.buf) > 0 {
			toCopy := len(r.buf)
			if toCopy > n-offset {
				toCopy = n - offset
			}
			copy(result[offset:], r.buf[:toCopy])
			r.buf = r.buf[toCopy:]
			offset += toCopy
			continue
		}
		data, err := r.transport.Recv()
		if err != nil {
			return nil, err
		}
		toCopy := len(data)
		if toCopy > n-offset {
			toCopy = n - offset
		}
		copy(result[offset:], data[:toCopy])
		r.buf = data[toCopy:]
		offset += toCopy
	}
	return result, nil
}
