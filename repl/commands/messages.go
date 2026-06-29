package commands

import (
	"fmt"

	"github.com/akovardin/gomax/repl/storage"
)

const CommandMessages = ".messages"

type Messages struct {
	storage *storage.Stotage
}

func NewMessages(storage *storage.Stotage) *Messages {
	return &Messages{
		storage: storage,
	}
}

func (s *Messages) Help() string {
	return ".messages - Выводит список новых сообщений"
}

func (s *Messages) Do(args []string) error {
	for _, m := range s.storage.List() {
		fmt.Printf("%s | %s\n\n", m.ChatID.String(), m.Text)
	}

	return nil
}
