package main

import (
	"time"

	"github.com/RussellLuo/timingwheel"
)

type Storage interface {
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) error
	Expire(key []byte, s float64) error
}

type MemoryStorage struct {
	data map[string][]byte
	tw   *timingwheel.TimingWheel
}

func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		data: make(map[string][]byte),
		tw:   timingwheel.NewTimingWheel(1*time.Millisecond, 1000),
	}
	storage.tw.Start()
	return storage
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

func (s *MemoryStorage) Expire(key []byte, seconds float64) error {
	s.tw.AfterFunc(time.Duration(seconds*float64(time.Second)), func() {
		s.Del(key)
	})
	return nil
}
