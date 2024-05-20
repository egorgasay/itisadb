package storage

import (
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (s *Storage) NewUser(user models.User) (r gost.ResultN) {
	s.users.Lock()
	defer s.users.Unlock()


	s.users.changeID++
	user.SetChangeID(s.users.changeID)
	s.users.Put(user.Login, user)

	return r.Ok()
}

func (s *Storage) GetUserByName(username string) (r gost.Result[models.User]) {
	s.users.RLock()
	defer s.users.RUnlock()

	user, ok := s.users.Get(username)
	if !ok || !user.Active {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(user)
}

func (s *Storage) DeleteUser(login string) (r gost.Result[bool]) {
	s.users.Lock()
	defer s.users.Unlock()

	val, ok := s.users.Get(login)
	if !ok || !val.Active {
		return r.Ok(false)
	}

	s.users.changeID++

	val.Active = false
	val.SetChangeID(s.users.changeID)
	s.users.Put(login, val)


	return r.Ok(true)
}

func (s *Storage) SaveUser(user models.User) (r gost.ResultN) {
	s.users.Lock()
	defer s.users.Unlock()
	
	if val, ok := s.users.Get(user.Login); !ok || !val.Active {
		return r.Err(constants.ErrNotFound)
	}

	s.users.changeID++

	s.users.Put(user.Login, user)
	user.SetChangeID(s.users.changeID)


	return r.Ok()
}

func (s *Storage) GetUserLevel(login string) (r gost.Result[models.Level]) {
	s.users.RLock()
	defer s.users.RUnlock()

	val, ok := s.users.Get(login)
	if !ok || !val.Active {
		return r.Err(constants.ErrNotFound)
	}

	return r.Ok(val.Level)
}

func (s *Storage) SetUserChangeID(id uint64) {
	s.users.Lock()
	defer s.users.Unlock()

	s.users.changeID = id
}
