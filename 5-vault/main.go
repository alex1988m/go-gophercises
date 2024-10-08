package main

import (
	"os"

	"github.com/alex1988m/go-gophercises/5-vault/commands"
	"github.com/alex1988m/go-gophercises/5-vault/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log *logrus.Logger
func init() {
	log = logger.NewLogger()

	if os.Getenv("CIPHER_KEY") == "" {
		if err := godotenv.Load(); err != nil {
			log.WithError(err).Warn("Error loading .env file")
		}
	}

	if os.Getenv("CIPHER_KEY") == "" {
		log.Fatal("CIPHER_KEY is not set in environment variables or .env file")
	}
}

func main() {
	app := &cli.App{
		Name:  "vault",
		Usage: "A simple key-value store with encryption",
		Commands: []*cli.Command{
			{
				Name:   "set",
				Usage:  "Set a value in the vault",
				Action: commands.SetCommand,
			},
			{
				Name:   "get",
				Usage:  "Get a value from the vault",
				Action: commands.GetCommand,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
