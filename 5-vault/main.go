package main

import (
	"fmt"

	"github.com/alex1988m/go-gophercises/5-vault/vault"
)


func main() {
    key := []byte("change this pass")
    vault, err := vault.New(key)
    if err != nil {
        panic(err.Error())
    }

    // Set a value
    err = vault.Set("example", []byte("exampleplaintext"))
    if err != nil {
        panic(err.Error())
    }

    // Get the value
    decrypted, err := vault.Get("example")
    if err != nil {
        panic(err.Error())
    }
    fmt.Printf("Decrypted: %s\n", decrypted)

    // Try to get a non-existent key
    _, err = vault.Get("nonexistent")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
