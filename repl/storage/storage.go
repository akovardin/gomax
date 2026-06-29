package storage

import "github.com/akovardin/gomax/types"

type Stotage struct {
	messages []*types.Message
}

func New() *Stotage {
	return &Stotage{
		messages: []*types.Message{},
	}
}

func (s *Stotage) Add(msg *types.Message) {
	s.messages = append(s.messages, msg)
}

func (s *Stotage) List() []*types.Message {
	return s.messages
}
