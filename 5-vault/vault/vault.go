package vault

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
    "sync"
)

type Vault struct {
    key   []byte
    block cipher.Block
    data  sync.Map
}

func New(key []byte) (*Vault, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return &Vault{key: key, block: block}, nil
}

func (v *Vault) Get(key string) ([]byte, error) {
    value, ok := v.data.Load(key)
    if !ok {
        return nil, fmt.Errorf("key not found")
    }

    ciphertext, ok := value.([]byte)
    if !ok {
        return nil, fmt.Errorf("invalid value type")
    }

    if len(ciphertext) < aes.BlockSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(v.block, iv)
    decrypted := make([]byte, len(ciphertext))
    stream.XORKeyStream(decrypted, ciphertext)

    return decrypted, nil
}

func (v *Vault) Set(key string, value []byte) error {
    ciphertext := make([]byte, aes.BlockSize+len(value))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return err
    }

    stream := cipher.NewCFBEncrypter(v.block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], value)

    v.data.Store(key, ciphertext)
    return nil
}
