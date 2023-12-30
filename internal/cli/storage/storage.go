package storage

import (
	"errors"
	"slices"
	"strings"
	"sync"
)

type actions []string
type user string

type Storage struct {
	sync.RWMutex
	RAMStorage map[user]actions
}

func New() *Storage {
	return &Storage{
		RAMStorage: make(map[user]actions, 10),
	}
}

func (s *Storage) SaveCommand(cookie string, command string) {
	s.Lock()
	defer s.Unlock()
	s.RAMStorage[user(cookie)] = append(s.RAMStorage[user(cookie)], command)
}

func (s *Storage) GetHistory(cookie string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.RAMStorage[user(cookie)]
	if !ok {
		return "", errors.New("empty history")
	}

	slices.Reverse(val)

	return strings.Join(val, "<br>"), nil
}
