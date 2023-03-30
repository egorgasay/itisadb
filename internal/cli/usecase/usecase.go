package usecase

import (
	"log"
	"strings"

	"github.com/egorgasay/grpc-storage/internal/cli/commands"
	"github.com/egorgasay/grpc-storage/internal/cli/config"
	"github.com/egorgasay/grpc-storage/pkg/api/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/egorgasay/grpc-storage/internal/cli/storage"
)

type UseCase struct {
	storage *storage.Storage
	conn    *grpc.ClientConn
}

func New(cfg *config.Config, storage *storage.Storage) *UseCase {
	conn, err := grpc.Dial(cfg.Balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return &UseCase{conn: conn, storage: storage}
}

func (uc *UseCase) ProcessQuery(cookie string, line string) (string, error) {
	uc.storage.SaveCommand(cookie, line)
	cmds := commands.New(balancer.NewBalancerClient(uc.conn))
	split := strings.Split(line, " ")

	return cmds.Do(commands.Action(strings.ToLower(split[0])), split[1:]...)
}

func (uc *UseCase) History(cookie string) (string, error) {
	return uc.storage.GetHistory(cookie)
}
