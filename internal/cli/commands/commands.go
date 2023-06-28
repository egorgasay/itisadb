package commands

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	attach     = "attach"
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
	case attach:
		if len(args) < 2 {
			return "", ErrWrongInput
		}
		dst := args[0]
		src := args[1]
		if err := c.attach(dst, src); err != nil {
			return "", err
		}
		return fmt.Sprintf("status: ok, attached %s to %s", src, dst), nil
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
	m, err := c.cl.IndexToJSON(context.Background(), &balancer.BalancerIndexToJSONRequest{Name: name})
	if err != nil {
		return "", err
	}

	return m.Index, nil
}

func (c *Commands) index(act, name, key, value string) (string, error) {
	switch act {
	case set:
		return c.setIndex(name, key, value)
	case get:
		return c.IndexToJSON(name, key)
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

func (c *Commands) IndexToJSON(name, key string) (string, error) {
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

func (c *Commands) attach(dst string, src string) error {
	_, err := c.cl.AttachToIndex(context.Background(), &balancer.BalancerAttachToIndexRequest{
		Dst: dst,
		Src: src,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return fmt.Errorf("index not found")
		}

		if st.Code() == codes.Unavailable {
			return fmt.Errorf("server not available")
		}
		return err
	}

	return nil
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
