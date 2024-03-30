package storage

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (s *Storage) NewUser(user models.User) (r gost.Result[int]) {
	s.users.RLock()
	defer s.users.RUnlock()

	id := -1

	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Login == user.Login {
			id = k
			return true
		}
		return false
	})

	if id != -1 {
		return r.Ok(id)
	}

	id = s.users.Count()
	user.ID = id
	s.users.Put(id, user)

	return r.Ok(id)
}

func (s *Storage) GetUserByID(id int) (r gost.Result[models.User]) {
	s.users.RLock()
	defer s.users.RUnlock()

	val, ok := s.users.Get(id)
	if !ok {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(val)
}

func (s *Storage) GetUserByName(username string) (r gost.Result[models.User]) {
	s.users.RLock()
	defer s.users.RUnlock()

	var find *models.User
	s.users.Iter(func(k int, v models.User) (stop bool) {
		if v.Login == username {
			find = &v
			return true
		}
		return false
	})

	if find == nil {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(*find)
}

func (s *Storage) DeleteUser(id int) (r gost.Result[bool]) {
	s.users.Lock()
	defer s.users.Unlock()

	if !s.users.Has(id) {
		return r.Ok(false)
	}

	s.users.Delete(id)

	return r.Ok(true)
}

func (s *Storage) SaveUser(user models.User) (r gost.ResultN) {
	s.users.Lock()
	defer s.users.Unlock()

	if !s.users.Has(user.ID) {
		return r.Err(constants.ErrNotFound)
	}

	s.users.Put(user.ID, user)

	return r.Ok()
}

func (s *Storage) GetUserLevel(id int) (r gost.Result[models.Level]) {
	s.users.RLock()
	defer s.users.RUnlock()

	val, ok := s.users.Get(id)
	if !ok {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(val.Level)
}

func (s *Storage) SetUserChangeID(id uint64) {
	s.users.Lock()
	defer s.users.Unlock()

	s.users.changeID = id
}
