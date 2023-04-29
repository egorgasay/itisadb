package commands

import (
	"context"
	"errors"
	"fmt"
	"itisadb/pkg/api/balancer"
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
	get  = "get"
	set  = "set"
	uset = "uset"
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
	case set, uset:
		if len(args) < 2 {
			return "", ErrWrongInput
		}

		var server int32 = 0
		if len(args) > 2 {
			num, err := strconv.Atoi(args[len(args)-1])
			if err == nil {
				server = int32(num)
				return c.set(args[0], strings.Join(args[1:len(args)-1], " "), server, uset == act)
			}
		}

		return c.set(args[0], strings.Join(args[1:], " "), server, uset == act)
	}

	return "", ErrUnknownCMD
}

func (c *Commands) get(key string, server int32) (string, error) {
	resp, err := c.cl.Get(context.Background(), &balancer.BalancerGetRequest{Key: key, Server: server})
	if err != nil {
		return "", err
	}

	if resp.Value == "" {
		return "", ErrEmpty
	}
	return resp.Value, nil
}

func (c *Commands) set(key, value string, server int32, uniques bool) (string, error) {
	response, err := c.cl.Set(context.Background(), &balancer.BalancerSetRequest{Key: key, Value: value, Server: server,
		Uniques: uniques})
	if err != nil {
		return "", err
	}

	resp := ""

	switch response.SavedTo {
	case 0:
		resp = fmt.Sprintf("status: ok, saved on server #%d", response.SavedTo)
	case dbOnly:
		resp = fmt.Sprintf("status: ok, saved on the database")
	case setToAll:
		resp = fmt.Sprintf("status: ok, saved on all servers")
	case AllAndDB:
		resp = fmt.Sprintf("status: ok, saved on all servers and in the database")
	default:
		resp = fmt.Sprintf("status: ok, saved on server #%d", response.SavedTo)
	}

	return resp, nil
}
