package storage

type Storage interface {
	Add(string, interface{}) error
	Remove(string) error
	FindOne(string) (interface{}, error)
}

func NewStorage(storageType string) (Storage, error) {
	switch storageType {
	case "memory":
		return NewMemoryStorage()
	default:
		return NewMemoryStorage()
	}
}
