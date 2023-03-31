package usecase

import (
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-storage/internal/cli/commands"
	"grpc-storage/internal/cli/config"
	"grpc-storage/pkg/api/balancer"

	"grpc-storage/internal/cli/storage"
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
