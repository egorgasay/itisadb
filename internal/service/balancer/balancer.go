package balancer

import (
	"context"
	"errors"
	"runtime"

	"github.com/egorgasay/gost"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/service/logic"
)

type Balancer struct {
	logger *zap.Logger

	servers  domains.Servers
	storage  domains.Storage
	tlogger  domains.TransactionLogger
	session  domains.Session
	security domains.SecurityService
	*logic.Logic

	cfg config.Config

	pool chan struct{} // TODO: ADD TO CONFIG

	objectServers gost.RwLock[map[string]int32]
	keyServers    gost.RwLock[map[string]int32]
}

func New(
	ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
	storage domains.Storage,
	tlogger domains.TransactionLogger,
	servers domains.Servers,
	session domains.Session,
	security domains.SecurityService,
	logic *logic.Logic,
) (*Balancer, error) {
	var err error

	if err != nil && !errors.Is(err, constants.ErrAlreadyExists) {
		return nil, err
	}

	
	return &Balancer{
		logger:        logger,
		servers:       servers,
		storage:       storage,
		tlogger:       tlogger,
		session:       session,
		cfg:           cfg,
		pool:          make(chan struct{}, 20_000*runtime.NumCPU()), // TODO: MOVE TO CONFIG
		objectServers: gost.NewRwLock(make(map[string]int32)),
		keyServers:    gost.NewRwLock(make(map[string]int32)),
		security:      security,
		Logic:         logic,
	}, nil
}
