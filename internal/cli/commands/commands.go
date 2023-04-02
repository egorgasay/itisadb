package commands

import (
	"context"
	"errors"
	"fmt"
	"grpc-storage/pkg/api/balancer"
	"log"
	"strconv"
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

		var server = -1
		if len(args) != 1 {
			num, err := strconv.Atoi(args[len(args)-1])
			if err == nil {
				server = num
			}
		}

		return c.get(args[0], int32(server))
	case set:
		if len(args) < 2 {
			return "", ErrWrongInput
		}
		return c.set(args[0], strings.Join(args[1:], " "))
	}

	return "", ErrUnknownCMD
}

func (c *Commands) get(key string, server int32) (string, error) {
	res, err := c.cl.Get(context.Background(), &balancer.BalancerGetRequest{Key: key, Server: server})
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
