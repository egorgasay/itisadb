package syncer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"itisadb/internal/domains"

	"go.uber.org/zap"
)

type Syncer struct {
	servers domains.Servers
	repo    domains.Storage
	logger  *zap.Logger
	f       *os.File
}

var syncerIsRunning = false

func NewSyncer(servers domains.Servers, logger *zap.Logger, repo domains.Storage) (domains.Syncer, error) {
	s := &Syncer{
		servers: servers,
		repo:    repo,
		logger:  logger,
	}

	f, err := os.OpenFile("sync", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("can't open sync file: %w", err)
	}

	s.f = f

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("can't read sync file: %w", err)
	}

	if len(b) != 0 {
		syncID, err := strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse sync id from file: %w", err)
		}

		s.repo.SetUserChangeID(syncID)
	}

	return s, nil
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
	ctx := context.Background()

	// Получение последнего идентификатора изменений на проверяемом сервере
	r := server.GetLastUserChangeID(ctx) 
	if r.IsErr() {
		return r.Error()
	}

	syncID := r.Unwrap()

	// Получение последнего идентификатора изменений
	currentSyncID := s.repo.GetUserChangeID()
	if syncID == currentSyncID { 
		return nil // Если нет изменений - алгоритм завершает работу
	}

	s.logger.Info("syncing", zap.Uint64("sync_id", syncID), zap.Uint64("current_sync_id", currentSyncID))

	// Получение списка обновленных пользователей
	rUsers := s.repo.GetUsersFromChangeID(syncID)
	if rUsers.IsErr() {
		return rUsers.Error()
	}

	// Обновление пользователей
	rSync := server.Sync(ctx, currentSyncID, rUsers.Unwrap())
	if rSync.IsErr() {
		return rSync.Error()
	}

	_, err := s.f.WriteAt([]byte(fmt.Sprint(currentSyncID)), 0) // Сохраняем идентификатор в файл
	if err != nil {
		return fmt.Errorf("can't write sync id to file: %w", err)
	}

	return nil
}
