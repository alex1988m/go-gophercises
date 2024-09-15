package vault
import (
	"github.com/alex1988m/go-gophercises/5-vault/logger"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

var log *logrus.Logger = logger.NewLogger()
type Vault struct {
	cryptor Cryptor
	file    *fileVault
}

func New(key []byte, filePath string) (*Vault, error) {
	fv, err := newFileVault(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file vault")
	}

	cryptor := NewAESCryptor(key)
	return &Vault{cryptor: cryptor, file: fv}, nil
}
func (v *Vault) Get(key string) ([]byte, error) {
	value, ok := v.file.get(key)
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

	err = v.file.set(key, encrypted)
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to store encrypted value")
		return errors.Wrap(err, "failed to store encrypted value")
	}

	log.WithField("key", key).Info("Successfully encrypted and stored value")
	return nil
}
