package main

import (
	"github.com/akovardin/gomax"
	"github.com/akovardin/gomax/repl/commands"
	"github.com/akovardin/gomax/repl/storage"
	"github.com/akovardin/gomax/repl/types"
)

func main() {
	cfg := gomax.DefaultExtraConfig()
	cfg.LogLevel = "error"

	client := gomax.NewWebClient(
		"session-web.db", // sessionName
		".",              // workDir
		cfg,              // extraConfig
		nil,              // qrProvider
		nil,              // passwordProvider
	)

	storage := storage.New()

	cmds := map[string]types.Command{
		commands.CommandSend:     commands.NewSend(client),
		commands.CommandClear:    commands.NewClear(),
		commands.CommandExit:     commands.NewExit(),
		commands.CommandChats:    commands.NewChats(client),
		commands.CommandMessages: commands.NewMessages(storage),
	}

	hlp := commands.NewHelp(cmds)

	cmds[commands.CommandHelp] = types.Command(hlp)

	repl := Repl{
		client:   client,
		commands: cmds,
		storage:  storage,
	}

	repl.Run()
}
