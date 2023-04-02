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
var ErrUnknownServer = errors.New("the value does not exist")

const (
	get = "get"
	set = "set"
)

const (
	_ = iota * -1
	dbOnly
	setToAll
	AllAndDB
)

func (c *Commands) Do(act Action, args ...string) (string, error) {
	switch act {
	case get:
		if len(args) < 1 {
			return "", ErrWrongInput
		}

		var server = 0
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

		var server int32 = 0
		if len(args) > 2 {
			num, err := strconv.Atoi(args[len(args)-1])
			if err == nil {
				server = int32(num)
				return c.set(args[0], strings.Join(args[1:len(args)-1], " "), server)
			}
		}

		return c.set(args[0], strings.Join(args[1:], " "), server)
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

func (c *Commands) set(key, value string, server int32) (string, error) {
	res, err := c.cl.Set(context.Background(), &balancer.BalancerSetRequest{Key: key, Value: value, Server: server})
	if err != nil {
		return "", err
	}

	log.Println(res.String())

	resp := ""

	switch server {
	case 0:
		resp = fmt.Sprintf("status: %s, saved on server #%d", res.Status, res.SavedTo)
	case dbOnly:
		resp = fmt.Sprintf("status: %s, saved on the database", res.Status)
	case setToAll:
		resp = fmt.Sprintf("status: %s, saved on all servers", res.Status)
	case AllAndDB:
		resp = fmt.Sprintf("status: %s, saved on all servers and in the database", res.Status)
	default:
		resp = fmt.Sprintf("status: %s, saved on server #%d", res.Status, res.SavedTo)
	}

	return resp, nil
}
