package commands

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "add a task to the list",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("add task: ", cCtx.Args().First())
			return nil
		},
	}
}

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "print task list",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("list task: ")
			return nil
		},
	}
}

func NewDoCommand() *cli.Command {
	return &cli.Command{
		Name:    "do",
		Aliases: []string{"d"},
		Usage:   "complete a task",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("do task: ", cCtx.Args().First())
			return errors.New("not implemented")
		},
	}
}
