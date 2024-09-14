package main

import (
 "github.com/sirupsen/logrus"
 
 "github.com/alex1988m/go-gophercises/5-vault/vault"
)

func init() {
    // Configure logrus
    logrus.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    logrus.SetReportCaller(true)
    logrus.SetLevel(logrus.DebugLevel)
}
func main() {
    key := []byte("change this pasd")
    vault := vault.New(key)

    // Set a value
    err := vault.Set("example", []byte("exampleplaintext"))
    if err != nil {
        logrus.WithError(err).Error("Failed to set value in vault")
    } else {
        logrus.Info("Successfully set value in vault")
    }

    // Get the value
    decrypted, err := vault.Get("example")
    if err != nil {
        logrus.WithError(err).Error("Failed to get value from vault")
    } else {
        logrus.WithField("decrypted", string(decrypted)).Info("Successfully retrieved and decrypted value")
    }

    // Try to get a non-existent key
    _, err = vault.Get("nonexistent")
    if err != nil {
        logrus.WithError(err).Warn("Attempted to get non-existent key")
    }
}
