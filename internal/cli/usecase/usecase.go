package usecase

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"itisadb/config"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"itisadb/internal/cli/commands"
	"itisadb/pkg/api"

	"itisadb/internal/cli/storage"
)

type UseCase struct {
	storage *storage.Storage
	conn    api.ItisaDBClient

	// sdk *itisadb.Client

	cmds   *commands.Commands
	tokens map[string]string
}

func New(cfg config.WebAppConfig, storage *storage.Storage, balancer string) *UseCase {
	conn, err := grpc.Dial(balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	b := api.NewItisaDBClient(conn)
	cmds := commands.New(b)

	return &UseCase{conn: b, storage: storage, cmds: cmds}
}

func (uc *UseCase) ProcessQuery(token string, line string) (string, error) {
	uc.storage.SaveCommand(token, line)
	split := strings.Split(line, " ")

	res, err := uc.cmds.Do(token, commands.Action(strings.ToLower(split[0])), split[1:]...)
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
			metadata.New(map[string]string{"token": token})), &api.ServersRequest{})
	return servers.ServersInfo, err
}

func (uc *UseCase) Authenticate(ctx context.Context, username, password string) (string, error) {
	resp, err := uc.conn.Authenticate(ctx, &api.AuthRequest{Login: username, Password: password})
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("failed to authenticate"))
	}

	// TODO: save token

	return resp.Token, nil
}
