package tcp

import (
	"encoding/binary"

	"github.com/akovardin/gomax/protocol"
)

const headerSize = 10

type TcpPacketFramer struct{}

func (f *TcpPacketFramer) Pack(ver, cmd, seq, opcode, flags int, payloadBytes []byte) []byte {
	payloadLen := len(payloadBytes)
	packedLen := uint32(((flags & 0xFF) << 24) | (payloadLen & 0x00FFFFFF))

	buf := make([]byte, headerSize+payloadLen)

	buf[0] = byte(ver)
	buf[1] = byte(cmd)
	binary.BigEndian.PutUint16(buf[2:4], uint16(seq))
	binary.BigEndian.PutUint16(buf[4:6], uint16(opcode))
	binary.BigEndian.PutUint32(buf[6:10], packedLen)

	copy(buf[headerSize:], payloadBytes)

	return buf
}

func (f *TcpPacketFramer) Unpack(data []byte) *protocol.PackedPacket {
	if len(data) < headerSize {
		return nil
	}

	packedLen := binary.BigEndian.Uint32(data[6:10])

	ver := int(data[0])
	cmd := int(data[1])
	seq := int(binary.BigEndian.Uint16(data[2:4]))
	opcode := int(binary.BigEndian.Uint16(data[4:6]))
	flags := int((packedLen >> 24) & 0xFF)
	payloadLen := int(packedLen & 0x00FFFFFF)

	header := protocol.TcpPacketHeader{
		Ver:        ver,
		Cmd:        cmd,
		Seq:        seq,
		Opcode:     opcode,
		Flags:      flags,
		PayloadLen: payloadLen,
	}

	payloadStart := headerSize
	payloadEnd := payloadStart + payloadLen
	if payloadEnd > len(data) {
		payloadEnd = len(data)
	}

	return &protocol.PackedPacket{
		Header:       header,
		PayloadBytes: data[payloadStart:payloadEnd],
	}
}

func (f *TcpPacketFramer) UnpackHeader(data []byte) (payloadLen int, ok bool) {
	if len(data) < headerSize {
		return 0, false
	}
	packedLen := binary.BigEndian.Uint32(data[6:10])
	payloadLen = int(packedLen & 0x00FFFFFF)
	return payloadLen, true
}
