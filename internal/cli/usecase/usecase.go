package usecase

import (
	"errors"
	"github.com/egorgasay/grpc-storage/internal/cli/commands"
	"github.com/egorgasay/grpc-storage/internal/cli/config"
	"github.com/egorgasay/grpc-storage/pkg/api/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

type UseCase struct {
	conn *grpc.ClientConn
}

func New(cfg *config.Config) *UseCase {
	conn, err := grpc.Dial(cfg.Balancer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return &UseCase{conn: conn}
}

func (uc *UseCase) ProcessQuery(cfg *config.Config, line string) (string, error) {
	if cfg == nil {
		return "", errors.New("empty config")
	}

	cmds := commands.New(balancer.NewBalancerClient(uc.conn))
	split := strings.Split(line, " ")

	return cmds.Do(commands.Action(strings.ToLower(split[0])), split[1:]...)
}
