package main

type Storage interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) error
}

type MemoryStorage struct {
	data map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (s *MemoryStorage) Set(key []byte, value []byte) error {
	s.data[string(key)] = value
	return nil
}

func (s *MemoryStorage) Get(key []byte) ([]byte, error) {
	value, ok := s.data[string(key)]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (s *MemoryStorage) Del(key []byte) error {
	delete(s.data, string(key))
	return nil
}
