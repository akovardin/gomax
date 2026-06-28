package main

import (
	"log"

	"github.com/akovardin/gomax"
	"github.com/akovardin/gomax/types"
)

func main() {
	extraConfig := gomax.DefaultExtraConfig()
	extraConfig.LogLevel = "debug"

	client := gomax.NewClient(
		"+79818463973",      // phone
		"session-mobile.db", // sessionName
		".",                 // workDir
		extraConfig,         // extraConfig
		nil,                 // smsCodeProvider
		nil,                 // passwordProvider
	)

	client.OnStart()(func(c interface{}) error {
		cl := c.(*gomax.Client)
		if me := cl.Me(); me != nil && me.Contact != nil {
			log.Printf("Запущен: %s", me.Contact.ID)
		}
		return nil
	})

	client.OnMessage()(func(event interface{}, c interface{}) error {
		msg := event.(*types.Message)
		log.Printf("chat=%v sender=%v text=%q id=%v", msg.ChatID, msg.Sender, msg.Text, msg.ID)
		return nil
	})

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
