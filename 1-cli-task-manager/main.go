package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alex1988m/go-gophercises/1-cli-task-manager/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	var err error
	app := &cli.App{
		Commands: []*cli.Command{
			commands.NewAddCommand(),
			commands.NewListCommand(),
			commands.NewDoCommand(),
		},
	}
	err = commands.InitStore()
	if err != nil {
		fmt.Println("failed to init store")
		log.Fatal(err)
	}

	defer func() {
		err = commands.CloseStore()
		if err != nil {
			fmt.Println("failed to close store")
			log.Fatal(err)
		}
	}()

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
