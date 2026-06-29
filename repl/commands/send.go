package commands

import (
	"errors"
	"log"
	"strconv"

	"github.com/akovardin/gomax"
)

const CommandSend = ".send"

type Send struct {
	client *gomax.WebClient
}

func NewSend(client *gomax.WebClient) *Send {
	return &Send{
		client: client,
	}
}

func (s *Send) Help() string {
	return ".send - Отправить сообщение"
}

func (s *Send) Do(args []string) error {
	if len(args) != 2 {
		return errors.New("нужно указать чат и сообщение")
	}

	chatID, _ := strconv.Atoi(args[0])
	message := args[1]

	msg, err := s.client.API().Messages.SendMessage(chatID, message, nil, nil, false)
	if err != nil {
		log.Printf("error on send message: %v, %v\n", err, msg)
	}

	return nil
}
