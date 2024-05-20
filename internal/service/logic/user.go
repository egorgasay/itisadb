package logic

import (
	"context"

	"itisadb/internal/constants"
	"itisadb/internal/models"

	"github.com/egorgasay/gost"
)

func (l *Logic) NewUser(ctx context.Context, claims gost.Option[models.UserClaims], user models.User) (r gost.ResultN) {
	if !l.security.HasPermission(claims, user.Level) {
		return r.Err(constants.ErrForbidden)
	}

	user.Active = true
	if rUser := l.storage.NewUser(user); r.IsErr() {
		return r.Err(rUser.Error())
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteNewUser(user)
	}

	return r.Ok()
}

func (l *Logic) DeleteUser(ctx context.Context, claims gost.Option[models.UserClaims], login string) (r gost.Result[bool]) {
	rUser := l.storage.GetUserByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	user := rUser.Unwrap()

	if !user.Active {
		return r.Err(constants.ErrNotFound)
	}

	if !l.security.HasPermission(claims, user.Level) {
		return r.Err(constants.ErrForbidden)
	}


	if r := l.storage.DeleteUser(user.Login); r.IsErr() {
		return r
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteDeleteUser(login)
	}

	return r
}

func (l *Logic) ChangePassword(ctx context.Context, claims gost.Option[models.UserClaims], login string, password string) (r gost.ResultN) {
	rUser := l.storage.GetUserByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	user := rUser.Unwrap()

	if !l.security.HasPermission(claims, user.Level) {
		return r.Err(constants.ErrForbidden)
	}

	user.Password = password

	if r := l.storage.SaveUser(user); r.IsErr() {
		return r
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteNewUser(user)
	}

	return r
}

func (l *Logic) ChangeLevel(ctx context.Context, claims gost.Option[models.UserClaims], login string, level models.Level) (r gost.ResultN) {
	rUser := l.storage.GetUserByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	user := rUser.Unwrap()

	if !l.security.HasPermission(claims, user.Level) {
		return r.Err(constants.ErrForbidden)
	}
	
	user.Level = level

	if r := l.storage.SaveUser(user); r.IsErr() {
		return r
	}

	if l.cfg.TransactionLogger.On {
		l.tlogger.WriteNewUser(user)
	}

	return r
}

func (l *Logic) GetLastUserChangeID(ctx context.Context) (r gost.Result[uint64]) {
	return r.Ok(l.storage.GetUserChangeID())
}

func (l *Logic) Sync(ctx context.Context, syncID uint64, users []models.User) (r gost.ResultN) {
	for _, user := range users {
		user.SetChangeID(syncID)
		l.storage.SaveUser(user)
	}

	l.storage.SetUserChangeID(syncID)

	return r.Ok()
}
