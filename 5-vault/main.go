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

	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}
}

func main() {
	storage, err := vault.NewFileStorage("vault.json")
	if err != nil {
		log.WithError(err).Fatal("Failed to create file storage")
	}

	key := []byte(os.Getenv("CIPHER_KEY"))
	if len(key) == 0 {
		log.Fatal("CIPHER_KEY environment variable is not set")
	}

	v, err := vault.New(key, storage)
	if err != nil {
		log.WithError(err).Fatal("Failed to create vault")
	}

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
