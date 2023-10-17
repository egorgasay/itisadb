package commands

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/pkg/api"
	"strconv"
	"strings"
)

type Commands struct {
	cl api.ItisaDBClient
}

func New(cl api.ItisaDBClient) *Commands {
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
	get    = "get"
	newEl  = "new"
	object = "object"
	show   = "show"
	attach = "attach"
)

const (
	level = "level"
	on    = "on"
)

const (
	_ = iota * -1
	dbOnly
	setToAll
	AllAndDB
)

func (c *Commands) Do(ctx context.Context, act Action, args ...string) (string, error) {
	switch strings.ToLower(string(act)) {
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

		return c.get(ctx, args[0], int32(server))
	case Set, uset:
		if len(args) < 2 {
			return "", ErrWrongInput
		}

		var server int32 = 0
		if len(args) > 2 {
			num, err := strconv.Atoi(args[len(args)-1])
			if err == nil {
				server = int32(num)
				return c.set(ctx, args[0], strings.Join(args[1:len(args)-1], " "), server, uset == act)
			}
		}

		return c.set(ctx, args[0], strings.Join(args[1:], " "), server, uset == act)
	case newEl:
		if len(args) < 1 {
			return "", ErrWrongInput
		}

		switch strings.ToLower(args[0]) {
		case object:
			if len(args) < 2 {
				return "", ErrWrongInput
			}
			name := args[1]
			return c.newObject(ctx, name)
		default:
			return c.set(ctx, args[0], "", 0, false)
		}
	case object:
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
		return c.object(ctx, act, name, key, value)
	case show:
		if len(args) < 1 {
			return "", ErrWrongInput
		}

		switch strings.ToLower(args[0]) {
		case object:
			if len(args) < 2 {
				return "", ErrWrongInput
			}
			name := args[1]
			return c.showObject(ctx, name)
		default:
			return "", ErrUnknownCMD
		}
	case attach:
		if len(args) < 2 {
			return "", ErrWrongInput
		}
		dst := args[0]
		src := args[1]
		if err := c.attach(ctx, dst, src); err != nil {
			return "", err
		}
		return fmt.Sprintf("status: ok, attached %s to %s", src, dst), nil
	}

	return "", ErrUnknownCMD
}

func (c *Commands) newObject(ctx context.Context, name string) (string, error) {
	_, err := c.cl.Object(ctx, &api.ObjectRequest{Name: name})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, object %s created", name), nil
}

func (c *Commands) showObject(ctx context.Context, name string) (string, error) {
	m, err := c.cl.ObjectToJSON(context.Background(), &api.ObjectToJSONRequest{Name: name})
	if err != nil {
		return "", err
	}

	return m.Object, nil
}

func (c *Commands) object(ctx context.Context, act, name, key, value string) (string, error) {
	switch act {
	case Set:
		return c.setObject(ctx, name, key, value)
	case get:
		return c.ObjectToJSON(ctx, name, key)
	default:
		return "", fmt.Errorf("unknown action")
	}

}

func (c *Commands) setObject(ctx context.Context, name, key, value string) (string, error) {
	r, err := c.cl.SetToObject(ctx, &api.SetToObjectRequest{
		Object: name,
		Key:    key,
		Value:  value,
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, saved in object %s, on server #%d", name, r.SavedTo), nil
}

func (c *Commands) ObjectToJSON(ctx context.Context, name, key string) (string, error) {
	r, err := c.cl.GetFromObject(ctx, &api.GetFromObjectRequest{
		Object: name,
		Key:    key,
	})

	if err != nil {
		return "", err
	}

	return r.Value, nil
}

func (c *Commands) get(сtx context.Context, key string, server int32) (string, error) {
	resp, err := c.cl.Get(сtx, &api.GetRequest{Key: key, Options: &api.GetRequest_Options{Server: &server}})
	if err != nil {
		return "", err
	}

	if resp.Value == "" {
		return "", ErrEmpty
	}
	return resp.Value, nil
}

func (c *Commands) attach(ctx context.Context, dst string, src string) error {
	_, err := c.cl.AttachToObject(ctx, &api.AttachToObjectRequest{
		Dst: dst,
		Src: src,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.NotFound {
			return fmt.Errorf("object not found")
		}

		if st.Code() == codes.Unavailable {
			return fmt.Errorf("server not available")
		}
		return err
	}

	return nil
}

func (c *Commands) set(ctx context.Context, key, value string, server int32, uniques bool) (string, error) {
	response, err := c.cl.Set(ctx, &api.SetRequest{
		Key: key, Value: value,
		Options: &api.SetRequest_Options{
			Server:  &server,
			Uniques: uniques,
		},
	})
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
	// return "", nil
}

type Command interface {
	Action() string
	Args() []string
	Server() int32
	Unique() bool
	ReadOnly() bool
	Level() uint8
}

func ParseCommand(ctx context.Context, text string) (Command, error) {
	split := strings.Split(text, " ")
	if len(split) < 1 {
		return nil, fmt.Errorf("no command detected")
	}

	switch cmd := strings.ToLower(split[0]); cmd {
	case Set, uset, rset, urset:
		return ParseSet[SetCommand](cmd, split[1:])
	}

	return nil, fmt.Errorf("unknown command")
}
