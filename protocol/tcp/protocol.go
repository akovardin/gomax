package tcp

import (
	"fmt"

	"github.com/akovardin/gomax/protocol"
)

type TcpProtocol struct {
	framer          *TcpPacketFramer
	serializer      *MsgpackPayloadCodec
	compression     *Lz4BlockCompression
	zstdCompression *ZstdCompression
	payloadDecoder  *TcpPayloadDecoder
}

func (p *TcpProtocol) GetFramer() *TcpPacketFramer {
	return p.framer
}

func NewTcpProtocol() *TcpProtocol {
	zstdComp, _ := NewZstdCompression()
	return &TcpProtocol{
		framer:          &TcpPacketFramer{},
		serializer:      &MsgpackPayloadCodec{},
		compression:     &Lz4BlockCompression{},
		zstdCompression: zstdComp,
		payloadDecoder:  NewTcpPayloadDecoder(zstdComp),
	}
}

func (p *TcpProtocol) Version() int {
	return 10
}

func (p *TcpProtocol) Encode(frame *protocol.OutboundFrame) ([]byte, error) {
	payloadBytes, err := p.serializer.Encode(frame.Payload)
	if err != nil {
		return nil, fmt.Errorf("tcp encode: %w", err)
	}

	packed := p.framer.Pack(frame.Ver, frame.Cmd, frame.Seq, frame.Opcode, 0, payloadBytes)
	return packed, nil
}

func (p *TcpProtocol) Decode(raw []byte) (*protocol.InboundFrame, error) {
	packed := p.framer.Unpack(raw)
	if packed == nil {
		return nil, fmt.Errorf("tcp decode: incomplete header")
	}

	payload, err := p.payloadDecoder.Decode(packed.PayloadBytes, packed.Header.Flags)
	if err != nil {
		return nil, fmt.Errorf("tcp decode: %w", err)
	}

	seq := packed.Header.Seq
	return &protocol.InboundFrame{
		Opcode:  packed.Header.Opcode,
		Cmd:     packed.Header.Cmd,
		Seq:     &seq,
		Payload: payload,
	}, nil
}
