package commands

import (
	"os"
	"os/exec"
)

const CommandClear = ".clear"

type Clear struct {
}

func NewClear() *Clear {
	return &Clear{}
}

func (s *Clear) Help() string {
	return ".clear - Очистить консоль"
}

func (s *Clear) Do(args []string) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}
