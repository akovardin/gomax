package commands

import (
	"fmt"
	"strconv"

	"github.com/akovardin/gomax"
)

const CommandChats = ".chats"

type Chats struct {
	client *gomax.WebClient
}

func NewChats(client *gomax.WebClient) *Chats {
	return &Chats{
		client: client,
	}
}

func (s *Chats) Help() string {
	return ".chats - Список чатов"
}

func (s *Chats) Do(args []string) error {
	chats, err := s.client.API().Chats.FetchChats(nil)
	if err != nil {
		return err
	}

	owner := s.client.Me()

	users := map[int]string{}
	ids := []int{}

	for _, chat := range chats {
		if chat.IsChannel() && chat.Title != nil {
			fmt.Printf("%s | %s\n", *chat.Title, chat.ID.String())
		}

		if chat.IsDialog() {
			for p := range chat.Participants {
				if p != owner.Contact.ID.String() {
					id, _ := strconv.Atoi(p)
					ids = append(ids, id)
					users[id] = chat.ID.String()
				}
			}
		}
	}

	participants, err := s.client.API().Users.GetUsers(ids)
	if err != nil {
		return err
	}

	for _, p := range participants {
		name := ""
		if len(p.Names) > 0 {
			name = p.Names[0].Name
		}

		fmt.Printf("%s | %s\n", name, users[p.ID.Int()])
	}

	return nil
}
