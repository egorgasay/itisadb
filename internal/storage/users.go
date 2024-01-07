package storage

import (
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (s *Storage) NewUser(user models.User) (id int, err error) {
	s.users.RLock()
	defer s.users.RUnlock()

	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Login == user.Login {
			id = k
			return true
		}
		return false
	})

	if id != 0 {
		return id, constants.ErrAlreadyExists
	}

	id = s.users.Count()
	user.ID = id
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
		if v.Login == username {
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

func (s *Storage) DeleteUser(id int) error {
	s.users.Lock()
	defer s.users.Unlock()

	if !s.users.Has(id) {
		return constants.ErrNotFound
	}

	s.users.Delete(id)

	return nil
}

func (s *Storage) SaveUser(id int, user models.User) error {
	s.users.Lock()
	defer s.users.Unlock()

	if !s.users.Has(id) {
		return constants.ErrNotFound
	}

	s.users.Put(id, user)

	return nil
}
