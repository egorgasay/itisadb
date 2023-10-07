package storage

import (
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (s *Storage) CreateUser(user models.User) error {
	s.users.RLock()
	defer s.users.RUnlock()

	if s.users.Has(user.Username) {
		return constants.ErrAlreadyExists
	}

	s.users.Put(user.Username, user)

	return nil
}

func (s *Storage) GetUser(username string) (models.User, error) {
	s.users.RLock()
	defer s.users.RUnlock()

	val, ok := s.users.Get(username)
	if !ok {
		return models.User{}, constants.ErrNotFound
	}

	return val, nil
}
