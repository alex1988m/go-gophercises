package vault

import (
	"github.com/alex1988m/go-gophercises/5-vault/logger"
	"github.com/sirupsen/logrus"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

var log *logrus.Logger = logger.NewLogger()

type Vault struct {
	key  []byte
	file *fileVault
}

func (v *Vault) getBlock() (cipher.Block, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		log.WithError(err).Error("Failed to create new cipher")
		return nil, errors.Wrap(err, "failed to create new cipher")
	}
	return block, nil
}
func New(key []byte, filePath string) (*Vault, error) {
	fv, err := newFileVault(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file vault")
	}

	return &Vault{key: key, file: fv}, nil
}
func (v *Vault) Get(key string) ([]byte, error) {
	value, ok := v.file.get(key)
	if !ok {
		err := errors.New("key not found")

		return nil, err
	}

	ciphertext := value

	if len(ciphertext) < aes.BlockSize {
		err := errors.New("ciphertext too short")
		log.WithField("key", key).WithError(err).Error("Stored ciphertext is too short")
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	block, err := v.getBlock()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to get cipher block")
		return nil, errors.Wrap(err, "failed to get cipher block for decryption")
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	stream.XORKeyStream(decrypted, ciphertext)

	log.WithField("key", key).Info("Successfully decrypted value")
	return decrypted, nil
}
func (v *Vault) Set(key string, value []byte) error {
	ciphertext := make([]byte, aes.BlockSize+len(value))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to generate IV")
		return errors.Wrap(err, "failed to generate IV for encryption")
	}
	block, err := v.getBlock()
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to get cipher block")
		return errors.Wrap(err, "failed to get cipher block for encryption")
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], value)

	err = v.file.set(key, ciphertext)
	if err != nil {
		log.WithField("key", key).WithError(err).Error("Failed to store encrypted value")
		return errors.Wrap(err, "failed to store encrypted value")
	}

	log.WithField("key", key).Info("Successfully encrypted and stored value")
	return nil
}
