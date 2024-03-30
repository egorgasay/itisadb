package syncer

import (
	"context"
	"time"

	"go.uber.org/zap"
	"itisadb/internal/domains"
)

type Syncer struct {
	servers domains.Servers
	repo    domains.Storage
	logger  *zap.Logger
}

var syncerIsRunning = false

func NewSyncer(servers domains.Servers, logger *zap.Logger, repo domains.Storage) domains.Syncer {
	s := &Syncer{
		servers: servers,
		repo:    repo,
		logger:  logger,
	}

	// TODO: SET SYNC ID FROM FILE TO STORAGE

	return s
}

func (s Syncer) Start() {
	if syncerIsRunning {
		panic("Syncer is already running")
	}

	syncerIsRunning = true
	defer func() { syncerIsRunning = false }()

	for {
		if err := s.servers.Iter(s.syncServer); err != nil {
			s.logger.Error("can't iter over servers", zap.Error(err))
		}

		time.Sleep(5 * time.Second) // TODO: make it configurable
	}
}

func (s Syncer) syncServer(server domains.Server) error {
	ctx := context.TODO()

	r := server.GetLastUserChangeID(ctx)
	if r.IsErr() {
		return r.Error()
	}

	syncID := r.Unwrap()

	currentSyncID := s.repo.GetUserChangeID()

	if syncID == currentSyncID {
		return nil
	}

	rUsers := s.repo.GetUsersFromChangeID(syncID)
	if rUsers.IsErr() {
		return rUsers.Error()
	}

	rSync := server.Sync(context.TODO(), currentSyncID, rUsers.Unwrap())
	if rSync.IsErr() {
		return rSync.Error()
	}

	// TODO: SAVE SYNC ID TO FILE

	return nil
}
