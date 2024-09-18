package storage

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
	fs := &FileStorage{
		filePath: filePath,
		data:     make(map[string][]byte),
	}

	err := fs.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "failed to load vault file")
	}

	return fs, nil
}

func (fs *FileStorage) load() error {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &fs.data)
}

func (fs *FileStorage) save() error {
	data, err := json.Marshal(fs.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal vault data")
	}

	return os.WriteFile(fs.filePath, data, 0600)
}

func (fs *FileStorage) Get(key string) ([]byte, bool) {
	value, ok := fs.data[key]
	return value, ok
}

func (fs *FileStorage) Set(key string, value []byte) error {
	fs.data[key] = value
	return fs.save()
}
