package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akovardin/gomax"
	"github.com/akovardin/gomax/repl/storage"
	t "github.com/akovardin/gomax/repl/types"
	"github.com/akovardin/gomax/types"
)

var name string = "gomax"

type Repl struct {
	client   *gomax.WebClient
	storage  *storage.Stotage
	commands map[string]t.Command
}

func (r *Repl) Run() {
	fmt.Printf("Запускаю: %s\n", name)

	r.client.OnStart()(func(c interface{}) error {
		r.loop()

		return nil
	})

	r.client.OnMessage()(func(event interface{}, c interface{}) error {
		r.storage.Add(event.(*types.Message))

		return nil
	})

	if err := r.client.Start(); err != nil {
		log.Fatal(err)
	}
}

func (r *Repl) loop() {
	reader := bufio.NewScanner(os.Stdin)
	promt()

	for reader.Scan() {
		if err := reader.Err(); err != nil {
			log.Printf("error on scan: %v", err)

			return
		}

		text, args := input(reader.Text())

		if text == "" {
			continue
		}

		r.handler(text, args)

		promt()
	}

	fmt.Println()
}

func promt() {
	fmt.Print(name, "> ")
}

func invalid(text string) {
	fmt.Println(text, ": command not found")
}

func (r *Repl) handler(text string, args []string) {
	cmd, ok := r.commands[text]
	if !ok {
		invalid(text)

		return
	}

	if err := cmd.Do(args); err != nil {
		fmt.Println("err: ", err.Error())
	}
}

func input(text string) (string, []string) {
	output := strings.TrimSpace(text)
	output = strings.ToLower(output)

	params := strings.Fields(output)
	cmd := params[0]
	var args []string
	if len(params) > 1 {
		args = params[1:]
	}

	return cmd, args
}
