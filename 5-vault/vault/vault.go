package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func logEntry() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"package": "vault",
	})
}


type Vault struct {
    key  []byte
    data sync.Map
}
func (v *Vault) getBlock() (cipher.Block, error) {
    block, err := aes.NewCipher(v.key)
    if err != nil {
        logEntry().WithError(err).Error("Failed to create new cipher")
        return nil, errors.Wrap(err, "failed to create new cipher")
    }
    return block, nil
}
func New(key []byte) *Vault {
    return &Vault{key: key}
}
func (v *Vault) Get(key string) ([]byte, error) {
    value, ok := v.data.Load(key)
    if !ok {
        err := errors.New("key not found")
        logEntry().WithField("key", key).WithError(err).Warn("Attempted to get non-existent key")
        return nil, err
    }

    ciphertext, ok := value.([]byte)
    if !ok {
        err := errors.New("invalid value type")
        logEntry().WithField("key", key).WithError(err).Error("Invalid value type stored in vault")
        return nil, err
    }

    if len(ciphertext) < aes.BlockSize {
        err := errors.New("ciphertext too short")
        logEntry().WithField("key", key).WithError(err).Error("Stored ciphertext is too short")
        return nil, err
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]
    block, err := v.getBlock()
    if err != nil {
        logEntry().WithField("key", key).WithError(err).Error("Failed to get cipher block")
        return nil, errors.Wrap(err, "failed to get cipher block for decryption")
    }
    stream := cipher.NewCFBDecrypter(block, iv)
    decrypted := make([]byte, len(ciphertext))
    stream.XORKeyStream(decrypted, ciphertext)

    logEntry().WithField("key", key).Info("Successfully decrypted value")
    return decrypted, nil
}
func (v *Vault) Set(key string, value []byte) error {
    ciphertext := make([]byte, aes.BlockSize+len(value))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        logEntry().WithField("key", key).WithError(err).Error("Failed to generate IV")
        return errors.Wrap(err, "failed to generate IV for encryption")
    }
    block, err := v.getBlock()
    if err != nil {
        logEntry().WithField("key", key).WithError(err).Error("Failed to get cipher block")
        return errors.Wrap(err, "failed to get cipher block for encryption")
    }
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], value)

    v.data.Store(key, ciphertext)
    logEntry().WithField("key", key).Info("Successfully encrypted and stored value")
    return nil
}
