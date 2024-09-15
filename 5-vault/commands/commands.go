package commands

import (
	"fmt"
	"os"

	"github.com/alex1988m/go-gophercises/5-vault/vault"
	"github.com/urfave/cli/v2"
)

func GetCommand(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("Usage: vault get <key>")
	}

	key := c.Args().Get(0)

	v, err := initVault()
	if err != nil {
		return err
	}

	value, err := v.Get(key)
	if err != nil {
		return fmt.Errorf("Failed to get value: %w", err)
	}

	fmt.Printf("Value for key '%s': %s\n", key, string(value))
	return nil
}

func SetCommand(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("Usage: vault set <key> <value>")
	}

	key := c.Args().Get(0)
	value := c.Args().Get(1)

	v, err := initVault()
	if err != nil {
		return err
	}

	err = v.Set(key, []byte(value))
	if err != nil {
		return fmt.Errorf("Failed to set value: %w", err)
	}

	fmt.Printf("Successfully set value for key '%s'\n", key)
	return nil
}

func initVault() (*vault.Vault, error) {
	key := []byte(os.Getenv("CIPHER_KEY"))
	if len(key) == 0 {
		return nil, fmt.Errorf("CIPHER_KEY environment variable is not set")
	}

	v, err := vault.New(key, "vault.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to create vault: %w", err)
	}

	return v, nil
}
