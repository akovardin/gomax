package commands

import (
	"fmt"

	"github.com/akovardin/gomax/repl/types"
)

const CommandHelp = ".help"

type Help struct {
	commands map[string]types.Command
}

func NewHelp(commands map[string]types.Command) *Help {
	return &Help{
		commands: commands,
	}
}

func (s *Help) Help() string {
	return ".help - Показать доступные команды"
}

func (s *Help) Do(args []string) error {
	fmt.Printf(
		"Привет! Доступные команды: \n\n",
	)

	for _, c := range s.commands {
		fmt.Println(c.Help())
	}

	return nil
}
