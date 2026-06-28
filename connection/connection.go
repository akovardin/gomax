package connection

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/akovardin/gomax/api/core"
	"github.com/akovardin/gomax/connection/readers"
	"github.com/akovardin/gomax/protocol"
	"github.com/akovardin/gomax/transport"
)

type ConnectionManager struct {
	reader    readers.Reader
	transport transport.Transport
	proto     protocol.Protocol
	requests  *PendingRequests

	OnEvent func(frame *protocol.InboundFrame)
	OnClose func(err error)

	isOpen         bool
	connectionLost bool
	closeReported  bool
	seq            int

	recvCtx    context.Context
	recvCancel context.CancelFunc
	eventWg    sync.WaitGroup
	closeMu    sync.Mutex
}

func NewConnectionManager(reader readers.Reader, transport transport.Transport, proto protocol.Protocol) *ConnectionManager {
	return &ConnectionManager{
		reader:    reader,
		transport: transport,
		proto:     proto,
		requests:  NewPendingRequests(),
	}
}

func (c *ConnectionManager) Open() error {
	if err := c.transport.Connect(); err != nil {
		return err
	}
	c.isOpen = true
	c.recvCtx, c.recvCancel = context.WithCancel(context.Background())
	c.eventWg.Add(1)
	go c.recvLoop()
	return nil
}

func (c *ConnectionManager) Close() error {
	if !c.isOpen {
		return nil
	}
	c.isOpen = false

	if c.recvCancel != nil {
		c.recvCancel()
	}
	c.requests.CancelAll(errors.New("connection closed"))
	c.eventWg.Wait()

	return c.transport.Close()
}

func (c *ConnectionManager) Fail(err error) error {
	c.connectionLost = true
	c.requests.CancelAll(err)

	c.closeMu.Lock()
	if !c.closeReported && c.OnClose != nil {
		c.closeReported = true
		c.closeMu.Unlock()
		c.OnClose(err)
	} else {
		c.closeMu.Unlock()
	}

	return c.transport.Close()
}

func (c *ConnectionManager) Send(frame *protocol.OutboundFrame) error {
	if !c.isOpen {
		return errors.New("connection not open")
	}
	data, err := c.proto.Encode(frame)
	if err != nil {
		return err
	}
	core.LogDebug("send frame opcode=%d cmd=%d seq=%d bytes=%d", frame.Opcode, frame.Cmd, frame.Seq, len(data))
	return c.transport.Send(data)
}

func (c *ConnectionManager) Request(frame *protocol.OutboundFrame, timeout time.Duration) (*protocol.InboundFrame, error) {
	frame.Seq = c.NextSeq()
	ch := c.requests.Create(frame.Seq)
	defer c.requests.Discard(frame.Seq)

	if err := c.Send(frame); err != nil {
		return nil, err
	}

	select {
	case result, ok := <-ch:
		if !ok {
			return nil, errors.New("request cancelled")
		}
		return result, nil
	case <-time.After(timeout):
		return nil, errors.New("request timed out")
	}
}

func (c *ConnectionManager) Version() int {
	return c.proto.Version()
}

func (c *ConnectionManager) WaitClosed() error {
	c.eventWg.Wait()
	return nil
}

func (c *ConnectionManager) NextSeq() int {
	seq := c.seq
	c.seq = (c.seq + 1) % 65536
	return seq
}

func (c *ConnectionManager) IsOpen() bool {
	return c.isOpen && !c.connectionLost
}

func (c *ConnectionManager) recvLoop() {
	defer c.eventWg.Done()

	for {
		select {
		case <-c.recvCtx.Done():
			return
		default:
		}

		data, err := c.reader.Read()
		if err != nil {
			core.LogDebug("recv error: %v", err)
			c.handleRecvError(err)
			return
		}

		frame, err := c.proto.Decode(data)
		if err != nil {
			c.handleRecvError(err)
			return
		}

		c.handleInbound(frame)
	}
}

func (c *ConnectionManager) handleInbound(frame *protocol.InboundFrame) {
	core.LogDebug("recv frame opcode=%d cmd=%d seq=%v", frame.Opcode, frame.Cmd, frame.Seq)
	if frame.Cmd == int(protocol.CommandResponse) || frame.Cmd == int(protocol.CommandError) {
		if frame.Seq != nil {
			if c.requests.Resolve(*frame.Seq, frame) {
				return
			}
		}
	}

	if c.OnEvent != nil {
		c.eventWg.Add(1)
		go func() {
			defer c.eventWg.Done()
			c.OnEvent(frame)
		}()
	}
}

func (c *ConnectionManager) handleRecvError(err error) {
	c.connectionLost = true
	c.requests.CancelAll(err)

	c.closeMu.Lock()
	if !c.closeReported && c.OnClose != nil {
		c.closeReported = true
		c.closeMu.Unlock()
		c.OnClose(err)
	} else {
		c.closeMu.Unlock()
	}
}
