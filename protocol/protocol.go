package protocol

type Protocol interface {
	Version() int
	Encode(frame *OutboundFrame) ([]byte, error)
	Decode(raw []byte) (*InboundFrame, error)
}
