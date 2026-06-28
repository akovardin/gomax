package transport

type Transport interface {
	Connect() error
	Close() error
	Send(data []byte) error
	Recv() ([]byte, error)
	Connected() bool
}
