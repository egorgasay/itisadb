package usecase

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"itisadb/internal/cli/commands"
	"itisadb/internal/cli/config"
	"itisadb/pkg/api/balancer"

	"itisadb/internal/cli/storage"
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

	res, err := cmds.Do(commands.Action(strings.ToLower(split[0])), split[1:]...)
	if err != nil {
		return res, err
	}

	return strings.Replace(strings.Replace(res, "\n", "<br/>", -1), "\t", "&emsp;", -1), nil
}

func (uc *UseCase) History(cookie string) (string, error) {
	return uc.storage.GetHistory(cookie)
}

func (uc *UseCase) Servers() (string, error) {
	b := balancer.NewBalancerClient(uc.conn)
	servers, err := b.Servers(context.TODO(), &balancer.BalancerServersRequest{})
	return servers.ServersInfo, err
}
