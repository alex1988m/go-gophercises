package main

import (
	"log"
	"os"

	"github.com/alex1988m/go-gophercises/1-cli-task-manager/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			commands.NewAddCommand(),
			commands.NewListCommand(),
			commands.NewDoCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
