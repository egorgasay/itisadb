package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/itisadb-go-sdk"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"itisadb/config"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"itisadb/internal/cli/commands"

	api "github.com/egorgasay/itisadb-shared-proto/go"

	"itisadb/internal/cli/storage"
)

type UseCase struct {
	storage *storage.Storage
	conn    api.ItisaDBClient

	sdk *itisadb.Client

	cmds   *commands.Commands
	tokens map[string]string
	logger *zap.Logger
}

func New(cfg config.WebAppConfig, storage *storage.Storage, balancer string, lg *zap.Logger) *UseCase {
	conn, err := grpc.Dial(balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	b := api.NewItisaDBClient(conn)

	r := itisadb.New(context.TODO(), balancer)
	if err := r.Error(); err != nil {
		log.Fatal(err)
	}

	cmds := commands.New(b, r.Unwrap())

	return &UseCase{conn: b, storage: storage, cmds: cmds, tokens: make(map[string]string), logger: lg}
}

func (uc *UseCase) ProcessQuery(ctx context.Context, token string, line string) (string, error) {
	uc.storage.SaveCommand(token, line)
	split := strings.Split(line, " ")

	res, err := uc.cmds.Do(withAuth(ctx, token), strings.ToLower(split[0]), split[1:]...)
	if err != nil {
		uc.logger.Warn(err.Error())
		return res, err
	}

	//cmd, err := commands.ParseCommand(ctx, line)
	//if err != nil {
	//	return "", err
	//}

	return strings.Replace(strings.Replace(res, "\n", "<br/>", -1), "\t", "&emsp;", -1), nil
}

func (uc *UseCase) SendCommand(ctx context.Context, cmd commands.Command) error {
	server := cmd.Server()
	readonly := cmd.Mode()
	level := cmd.Level()

	switch cmd.Action() {
	case commands.Set:
		//uc.sdk.SetOne(ctx, cmd.Args()[0], cmd.Args()[1], itisadb.SetOptions{
		//	Server:   &server,
		//	ReadOnly: readonly == 1,
		//})
		resp, err := uc.conn.Set(ctx, &api.SetRequest{
			Key:   cmd.Args()[0],
			Value: cmd.Args()[1],
			Options: &api.SetRequest_Options{
				Server:   &server,
				ReadOnly: readonly == 1, // TODO: fix
				Level:    api.Level(level),
			},
		})
		if err != nil {
			return err
		}
		_ = resp
		return nil
	default:
		return errors.New("unknown command")
	}
}

func (uc *UseCase) History(cookie string) (string, error) {
	return uc.storage.GetHistory(cookie)
}

func withAuth(ctx context.Context, token string) context.Context {
	return metadata.NewOutgoingContext(ctx,
		metadata.New(map[string]string{"token": token}))
}

func (uc *UseCase) Servers(ctx context.Context, token string) (string, error) {
	servers, err := uc.conn.Servers(withAuth(ctx, token), &api.ServersRequest{})
	return servers.ServersInfo, err
}

func (uc *UseCase) Authenticate(ctx context.Context, username, password string) (string, error) {
	resp, err := uc.conn.Authenticate(ctx, &api.AuthRequest{Login: username, Password: password})
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("failed to authenticate"))
	}

	// TODO: save token

	uc.tokens[username] = resp.Token

	return resp.Token, nil
}
