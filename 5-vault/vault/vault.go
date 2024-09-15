package vault

import (
	"github.com/alex1988m/go-gophercises/5-vault/logger"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

var log *logrus.Logger = logger.NewLogger()

type Vault struct {
	cryptor Cryptor
	storage Storage
}

func New(key []byte, storage Storage) (*Vault, error) {
	if storage == nil {
		return nil, errors.New("storage cannot be nil")
	}

	cryptor := NewAESCryptor(key)
	return &Vault{cryptor: cryptor, storage: storage}, nil
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
