package logic

import (
	"context"

	"github.com/egorgasay/gost"
	"itisadb/internal/models"
)

func (l *Logic) NewUser(ctx context.Context, _ gost.Option[models.UserClaims], user models.User) (r gost.ResultN) {
	if rUser := l.storage.NewUser(user); r.IsErr() {
		return r.Err(rUser.Error())
	}

	return r.Ok()
}

func (l *Logic) DeleteUser(ctx context.Context, _ gost.Option[models.UserClaims], login string) (r gost.Result[bool]) {
	rUser := l.storage.GetUserIDByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	return l.storage.DeleteUser(rUser.Unwrap())
}

func (l *Logic) ChangePassword(ctx context.Context, _ gost.Option[models.UserClaims], login string, password string) (r gost.ResultN) {
	rUser := l.storage.GetUserByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	user := rUser.Unwrap()
	user.Password = password

	return l.storage.SaveUser(user)
}

func (l *Logic) ChangeLevel(ctx context.Context, _ gost.Option[models.UserClaims], login string, level models.Level) (r gost.ResultN) {
	rUser := l.storage.GetUserByName(login)
	if rUser.IsErr() {
		return r.Err(rUser.Error())
	}

	user := rUser.Unwrap()
	user.Level = level

	return l.storage.SaveUser(user)
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
