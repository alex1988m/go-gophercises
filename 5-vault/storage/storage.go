package storage

// Storage defines the interface for storing and retrieving encrypted data
type Storage interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte) error
}
