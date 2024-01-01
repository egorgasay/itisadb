package balancer

import (
	"itisadb/internal/constants"
	"itisadb/internal/service/usecase"
)

type LocalServer struct {
	*usecase.UseCase
}

func NewLocalServer(uc *usecase.UseCase) *LocalServer {
	return &LocalServer{
		UseCase: uc,
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
