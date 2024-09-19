package vault

import (
	"github.com/alex1988m/go-gophercises/5-vault/cryptor"

	"github.com/alex1988m/go-gophercises/5-vault/logger"
	"github.com/alex1988m/go-gophercises/5-vault/storage"
	"github.com/sirupsen/logrus"

	"fmt"
	"os"

	"path/filepath"

	"github.com/pkg/errors"
)

var log *logrus.Logger = logger.NewLogger()
func NewVault() (*Vault, error) {
	key := []byte(os.Getenv("CIPHER_KEY"))
	if len(key) == 0 {
		return nil, fmt.Errorf("CIPHER_KEY environment variable is not set")
	}

	var vaultPath string
	localVaultPath := "vault.json"

	// Check if vault.json exists in the current directory
	if _, err := os.Stat(localVaultPath); err == nil {
		vaultPath = localVaultPath
	} else {
		// If not, use the home directory
		dir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Failed to get user home directory: %w", err)
		}
		vaultPath = filepath.Join(dir, "vault.json")
	}

	storage, err := storage.NewFileStorage(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file storage: %w", err)
	}

	cryptor := cryptor.NewAESCryptor(key)

	return &Vault{cryptor: cryptor, storage: storage}, nil
}

type Vault struct {
	cryptor cryptor.Cryptor
	storage storage.Storage
}

func (v *Vault) Get(key string) ([]byte, error) {
	value, ok := v.storage.Get(key)
	if !ok {
		err := errors.New("key not found")
		return nil, err
	}

	decrypted, err := v.cryptor.Decrypt(value)
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to decrypt value")
		return nil, errors.Wrap(err, "failed to decrypt value")
	}

	log.WithField("key", key).Info("Successfully decrypted value")
	return decrypted, nil
}

func (v *Vault) Set(key string, value []byte) error {
	encrypted, err := v.cryptor.Encrypt(value)
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to encrypt value")
		return errors.Wrap(err, "failed to encrypt value")
	}

	err = v.storage.Set(key, encrypted)
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to store encrypted value")
		return errors.Wrap(err, "failed to store encrypted value")
	}

	log.WithField("key", key).Info("Successfully encrypted and stored value")
	return nil
}
