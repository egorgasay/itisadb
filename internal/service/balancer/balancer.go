package balancer

import (
	"context"
	"go.uber.org/zap"
	"itisadb/config"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

type Balancer struct {
	logger *zap.Logger

	servers domains.Servers
	storage domains.Storage
	tlogger domains.TransactionLogger
	session domains.Session

	cfg config.Config
}

func (b *Balancer) Set(ctx context.Context, userID int, key, val string, opts models.SetOptions) (int32, error) {

}
