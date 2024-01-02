package balancer

import (
	"context"
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/models"
	"itisadb/internal/service/logic"
	"itisadb/pkg"
)

type LocalServer struct {
	*logic.Logic
	ram gost.RwLock[models.RAM]
}

func NewLocalServer(uc *logic.Logic) *LocalServer {
	return &LocalServer{
		Logic: uc,
		ram:   gost.NewRwLock(models.RAM{}),
	}
}

func (s *LocalServer) Tries() uint32 {
	return 0
}

func (s *LocalServer) IncTries() {}

func (s *LocalServer) ResetTries() {}

func (s *LocalServer) Number() int32 {
	return constants.MainStorageNumber
}

func (s *LocalServer) RAM() models.RAM {
	r := s.ram.RBorrow()
	defer s.ram.Release()

	return r.Read()
}

func (s *LocalServer) RefreshRAM(_ context.Context) (res gost.Result[gost.Nothing]) {
	r := pkg.CalcRAM()
	if r.IsErr() {
		return res.Err(r.Error())
	}

	s.ram.WBorrow()
	s.ram.WReturn(r.Unwrap())

	return res.Ok(gost.Nothing{})
}
