package commands

import (
	"os"
)

const CommandExit = ".exit"

type Exit struct {
}

func NewExit() *Exit {
	return &Exit{}
}

func (s *Exit) Help() string {
	return ".exit - Выйти"
}

func (s *Exit) Do(args []string) error {
	os.Exit(0)

	return nil
}
