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

func NewSyncer(servers domains.Servers, logger *zap.Logger, repo domains.Storage) *Syncer {
	return &Syncer{
		servers: servers,
		repo:    repo,
		logger:  logger,
	}
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
	r := server.GetLastSyncID()
	if r.IsErr() {
		return r.Error()
	}

	syncID := r.Unwrap()

	currentSyncID := s.repo.GetCurrentSyncID()

	if syncID != currentSyncID {
		return nil
	}

	rUsers := s.repo.GetUsersFromSyncID(syncID)
	if rUsers.IsErr() {
		return rUsers.Error()
	}

	rSync := server.Sync(context.TODO(), currentSyncID, users.Unwrap())
	if rSync.IsErr() {
		return rSync.Error()
	}

	return nil
}
