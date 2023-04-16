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
	gc balancer.Balancer_GetClient
	sc balancer.Balancer_SetClient
}

func New(cl balancer.BalancerClient) *Commands {
	gc, err := cl.Get(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	sc, err := cl.Set(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	return &Commands{
		cl: cl,
		gc: gc,
		sc: sc,
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
	err := c.gc.Send(&balancer.BalancerGetRequest{Key: key, Server: server})
	if err != nil {
		return "", err
	}

	recv, err := c.gc.Recv()
	if err != nil {
		return "", err
	}

	if recv.Value == "" {
		return "", ErrEmpty
	}

	log.Println(recv.Value)
	return recv.Value, nil
}

func (c *Commands) set(key, value string, server int32) (string, error) {
	err := c.sc.Send(&balancer.BalancerSetRequest{Key: key, Value: value, Server: server})
	if err != nil {
		return "", err
	}

	recv, err := c.sc.Recv()
	if err != nil {
		return "", err
	}

	log.Println(recv.String())

	resp := ""

	switch server {
	case 0:
		resp = fmt.Sprintf("status: %s, saved on server #%d", recv.Status, recv.SavedTo)
	case dbOnly:
		resp = fmt.Sprintf("status: %s, saved on the database", recv.Status)
	case setToAll:
		resp = fmt.Sprintf("status: %s, saved on all servers", recv.Status)
	case AllAndDB:
		resp = fmt.Sprintf("status: %s, saved on all servers and in the database", recv.Status)
	default:
		resp = fmt.Sprintf("status: %s, saved on server #%d", recv.Status, recv.SavedTo)
	}

	return resp, nil
}
