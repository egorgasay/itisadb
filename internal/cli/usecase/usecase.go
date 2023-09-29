package usecase

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
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
	conn    balancer.BalancerClient

	// sdk *itisadb.Client

	cmds   *commands.Commands
	tokens map[string]string
}

func New(cfg *config.Config, storage *storage.Storage) *UseCase {
	conn, err := grpc.Dial(cfg.Balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	b := balancer.NewBalancerClient(conn)
	cmds := commands.New(b)

	return &UseCase{conn: b, storage: storage, cmds: cmds}
}

func (uc *UseCase) ProcessQuery(cookie string, line string) (string, error) {
	uc.storage.SaveCommand(cookie, line)
	split := strings.Split(line, " ")

	res, err := uc.cmds.Do(commands.Action(strings.ToLower(split[0])), split[1:]...)
	if err != nil {
		return res, err
	}

	return strings.Replace(strings.Replace(res, "\n", "<br/>", -1), "\t", "&emsp;", -1), nil
}

func (uc *UseCase) History(cookie string) (string, error) {
	return uc.storage.GetHistory(cookie)
}

func (uc *UseCase) Servers(ctx context.Context, token string) (string, error) {
	servers, err := uc.conn.Servers(
		metadata.NewOutgoingContext(ctx,
			metadata.New(map[string]string{"token": token})), &balancer.BalancerServersRequest{})
	return servers.ServersInfo, err
}

func (uc *UseCase) Authenticate(ctx context.Context, username, password string) (string, error) {
	resp, err := uc.conn.Authenticate(ctx, &balancer.BalancerAuthRequest{Login: username, Password: password})
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("failed to authenticate"))
	}

	// TODO: save token

	return resp.Token, nil
}
