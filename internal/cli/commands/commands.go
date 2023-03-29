package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/grpc-storage/pkg/api/balancer"
	"log"
	"strings"
)

type Commands struct {
	cl balancer.BalancerClient
}

func New(cl balancer.BalancerClient) *Commands {
	return &Commands{
		cl: cl,
	}
}

type Action string

var ErrWrongInput = errors.New("wrong input")
var ErrUnknownCMD = errors.New("unknown cmd")
var ErrEmpty = errors.New("the value does not exist")

const (
	get = "get"
	set = "set"
)

func (c *Commands) Do(act Action, args ...string) (string, error) {
	switch act {
	case get:
		if len(args) < 1 {
			return "", ErrWrongInput
		}
		return c.get(args[0])
	case set:
		if len(args) < 2 {
			return "", ErrWrongInput
		}
		return c.set(args[0], strings.Join(args[1:], " "))
	}

	return "", ErrUnknownCMD
}

func (c *Commands) get(key string) (string, error) {
	res, err := c.cl.Get(context.Background(), &balancer.BalancerGetRequest{Key: key})
	if err != nil {
		return "", err
	} else if res.Value == "" {
		return "", ErrEmpty
	}

	log.Println(res.Value)
	return res.Value, nil
}

func (c *Commands) set(key, value string) (string, error) {
	res, err := c.cl.Set(context.Background(), &balancer.BalancerSetRequest{Key: key, Value: value})
	if err != nil {
		return "", err
	}

	log.Println(res.String())
	return fmt.Sprintf("status: %s, saved to server #%d", res.Status, res.SavedTo), nil
}
