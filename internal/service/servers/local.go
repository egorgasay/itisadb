package servers

import (
	"context"

	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
	"itisadb/internal/service/logic"
	"itisadb/pkg"

	"github.com/egorgasay/gost"
)

type LocalServer struct {
	*logic.Logic
	storage domains.Storage
	config  config.Config
	ram     gost.RwLock[models.RAM]
}

func NewLocalServer(uc *logic.Logic) *LocalServer {
	return &LocalServer{
		Logic: uc,
		ram:   gost.NewRwLock(models.RAM{}),
	}
}

func (s *LocalServer) IsOffline() bool { return false }

func (s *LocalServer) Reconnect(_ context.Context) (res gost.ResultN) {
	return
}

func (s *LocalServer) ResetTries() {}

func (s *LocalServer) Number() int32 {
	return constants.LocalServerNumber
}

func (s *LocalServer) RAM() models.RAM {
	r := s.ram.RBorrow()
	defer s.ram.Release()

	return r.Read()
}

func (s *LocalServer) RefreshRAM(_ context.Context) (res gost.ResultN) {
	r := pkg.CalcRAM()
	if r.IsErr() {
		return res.Err(r.Error())
	}

	s.ram.SetWithLock(r.Unwrap())

	return res.Ok()
}

func (s *LocalServer) NewUser(ctx context.Context, claims gost.Option[models.UserClaims], user models.User) (r gost.ResultN) {
	if s.config.Balancer.On {
		return r.Ok()
	}

	return s.Logic.NewUser(ctx, claims, user)
}

func (s *LocalServer) DeleteUser(ctx context.Context, claims gost.Option[models.UserClaims], login string) (r gost.Result[bool]) {
	if s.config.Balancer.On {
		return r.Ok(false)
	}

	return s.Logic.DeleteUser(ctx, claims, login)
}

func (s *LocalServer) ChangePassword(ctx context.Context, claims gost.Option[models.UserClaims], login string, password string) (r gost.ResultN) {
	if s.config.Balancer.On {
		return r.Ok()
	}

	return s.Logic.ChangePassword(ctx, claims, login, password)
}

func (s *LocalServer) ChangeLevel(ctx context.Context, claims gost.Option[models.UserClaims], login string, level models.Level) (r gost.ResultN) {
	if s.config.Balancer.On {
		return r.Ok()
	}

	return s.Logic.ChangeLevel(ctx, claims, login, level)
}

func (s *LocalServer) Address() string {
	return ""
}
