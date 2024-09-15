package vault

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type FileStorage struct {
	filePath string
	data     map[string][]byte
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fv := &fileVault{
		filePath: filePath,
		data:     make(map[string][]byte),
	}

	err := fv.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "failed to load vault file")
	}

	return fv, nil
}

func (fv *fileVault) load() error {
	data, err := os.ReadFile(fv.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &fv.data)
}

func (fv *fileVault) save() error {
	data, err := json.Marshal(fv.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal vault data")
	}

	return os.WriteFile(fv.filePath, data, 0600)
}

func (fs *FileStorage) Get(key string) ([]byte, bool) {
	value, ok := fv.data[key]
	return value, ok
}

func (fs *FileStorage) Set(key string, value []byte) error {
	fv.data[key] = value
	return fv.save()
}
