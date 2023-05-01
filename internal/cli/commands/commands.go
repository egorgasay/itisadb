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
	get        = "get"
	set        = "set"
	uset       = "uset"
	new_index  = "new_index"
	index      = "index"
	show_index = "show_index"
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
	case new_index:
		if len(args) < 1 {
			return "", ErrWrongInput
		}
		name := args[0]
		return c.newIndex(name)
	case index:
		if len(args) < 3 {
			return "", ErrWrongInput
		}
		name := args[0]
		act := args[1]
		key := args[2]
		value := ""
		if len(args) > 3 {
			value = strings.Join(args[3:], " ")
		}
		return c.index(act, name, key, value)
	case show_index:
		if len(args) < 1 {
			return "", ErrWrongInput
		}
		name := args[0]
		return c.showIndex(name)
	}

	return "", ErrUnknownCMD
}

func (c *Commands) newIndex(name string) (string, error) {
	_, err := c.cl.Index(context.Background(), &balancer.BalancerIndexRequest{Name: name})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, index %s created", name), nil
}

func (c *Commands) showIndex(name string) (string, error) {
	m, err := c.cl.GetIndex(context.Background(), &balancer.BalancerGetIndexRequest{Name: name})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("index %s:<br>", name))
	for k, v := range m.Index {
		sb.WriteString(fmt.Sprintf("%s: %s<br>", k, v))
	}

	return sb.String(), nil
}

func (c *Commands) index(act, name, key, value string) (string, error) {
	switch act {
	case set:
		return c.setIndex(name, key, value)
	case get:
		return c.getIndex(name, key)
	default:
		return "", fmt.Errorf("unknown action")
	}

}

func (c *Commands) setIndex(name, key, value string) (string, error) {
	r, err := c.cl.SetToIndex(context.Background(), &balancer.BalancerSetToIndexRequest{
		Index: name,
		Key:   key,
		Value: value,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, saved in index %s, on server #%d", name, r.SavedTo), nil
}

func (c *Commands) getIndex(name, key string) (string, error) {
	r, err := c.cl.GetFromIndex(context.Background(), &balancer.BalancerGetFromIndexRequest{
		Index: name,
		Key:   key,
	})

	if err != nil {
		return "", err
	}

	return r.Value, nil
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
