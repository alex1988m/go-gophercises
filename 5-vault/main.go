package main

import (
	
	"github.com/alex1988m/go-gophercises/5-vault/logger"	
	
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/alex1988m/go-gophercises/5-vault/vault"
)
var log *logrus.Logger

func init() {

	log = logger.NewLogger()


	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}
}

func main() {
	key := []byte(os.Getenv("CIPHER_KEY"))
	if len(key) == 0 {
		log.Fatal("CIPHER_KEY environment variable is not set")
	}

	v, err := vault.New(key, "vault.json")
	if err != nil {
		log.WithError(err).Fatal("Failed to create vault")
	}


	if err := v.Set("example", []byte("exampleplaintext")); err != nil {
		log.WithError(err).Error("Failed to set value in vault")
	} else {
		log.Info("Successfully set value in vault")
	}


	decrypted, err := v.Get("example")
	if err != nil {
		log.WithError(err).Error("Failed to get value from vault")
	} else {
		log.WithField("decrypted", string(decrypted)).Info("Successfully retrieved and decrypted value")
	}


	if _, err := v.Get("nonexistent"); err != nil {
		log.WithError(err).Warn("Attempted to get non-existent key")
	}
}
