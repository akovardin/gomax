package connection

import (
	"sync"

	"github.com/akovardin/gomax/protocol"
)

type PendingRequests struct {
	mu    sync.Mutex
	items map[int]chan *protocol.InboundFrame
}

func NewPendingRequests() *PendingRequests {
	return &PendingRequests{
		items: make(map[int]chan *protocol.InboundFrame),
	}
}

func (p *PendingRequests) Create(seq int) chan *protocol.InboundFrame {
	p.mu.Lock()
	defer p.mu.Unlock()
	ch := make(chan *protocol.InboundFrame, 1)
	p.items[seq] = ch
	return ch
}

func (p *PendingRequests) Resolve(seq int, frame *protocol.InboundFrame) bool {
	p.mu.Lock()
	ch, ok := p.items[seq]
	if ok {
		delete(p.items, seq)
	}
	p.mu.Unlock()
	if ok {
		ch <- frame
		return true
	}
	return false
}

func (p *PendingRequests) Discard(seq int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.items, seq)
}

func (p *PendingRequests) CancelAll(err error) {
	p.CancelAllWithError(err)
}

func (p *PendingRequests) CancelAllWithError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for seq, ch := range p.items {
		errorFrame := &protocol.InboundFrame{
			Error: err.Error(),
		}
		select {
		case ch <- errorFrame:
		default:
		}
		close(ch)
		delete(p.items, seq)
	}
}
