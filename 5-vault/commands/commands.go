package commands

import (
	"fmt"

	"github.com/alex1988m/go-gophercises/5-vault/vault"
	"github.com/urfave/cli/v2"
)

func GetCommand(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("Usage: vault get <key>")
	}

	key := c.Args().Get(0)
	v, err := vault.NewVault()
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

	v, err := vault.NewVault()
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
