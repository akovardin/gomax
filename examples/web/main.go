package main

import (
	"log"

	"github.com/akovardin/gomax"
	"github.com/akovardin/gomax/types"
)

func main() {
	cfg := gomax.DefaultExtraConfig()
	cfg.LogLevel = "debug"

	client := gomax.NewWebClient(
		"session-web.db", // sessionName
		".",              // workDir
		cfg,              // extraConfig
		nil,              // qrProvider
		nil,              // passwordProvider
	)

	client.OnStart()(func(c interface{}) error {
		log.Println("Web-клиент запущен")

		// отправка сообщения самому себе
		msg, err := client.API().Messages.SendMessage(0, "Web-клиент запущен", nil, nil, false)
		if err != nil {
			log.Printf("error on send message: %v, %v\n", err, msg)
		}

		log.Printf("sent message: %v\n", msg)

		return nil
	})

	client.OnMessage()(func(event interface{}, c interface{}) error {
		msg := event.(*types.Message)
		var chatID, sender types.FlexInt
		if msg.ChatID != nil {
			chatID = *msg.ChatID
		}
		if msg.Sender != nil {
			sender = *msg.Sender
		}
		log.Printf("%d %d %s | %v\n", chatID, sender, msg.Text, msg)

		return nil
	})

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
