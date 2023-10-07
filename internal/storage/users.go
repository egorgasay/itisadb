package storage

import (
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (s *Storage) CreateUser(user models.User) (id int, err error) {
	s.users.RLock()
	defer s.users.RUnlock()

	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Username == user.Username {
			id = k
			return true
		}
		return false
	})

	if id != 0 {
		return id, constants.ErrAlreadyExists
	}

	id = s.users.Count()
	s.users.Put(id, user)

	return id, nil
}

func (s *Storage) GetUserByID(id int) (models.User, error) {
	s.users.RLock()
	defer s.users.RUnlock()

	val, ok := s.users.Get(id)
	if !ok {
		return models.User{}, constants.ErrNotFound
	}

	return val, nil
}

func (s *Storage) GetUserByName(username string) (id int, u models.User, err error) {
	s.users.RLock()
	defer s.users.RUnlock()

	find := false
	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Username == username {
			id = k
			u = v
			find = true
			return true
		}
		return false
	})

	if !find {
		return id, u, constants.ErrNotFound
	}

	return id, u, nil
}
