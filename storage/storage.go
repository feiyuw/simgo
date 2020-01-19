package storage

type Storage interface {
	Add(interface{}, interface{}) error
	Remove(interface{}) error
	FindOne(interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
}

func NewStorage(storageType string) (Storage, error) {
	switch storageType {
	case "memory":
		return NewMemoryStorage()
	default:
		return NewMemoryStorage()
	}
}
