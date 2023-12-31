package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	api "github.com/egorgasay/itisadb-shared-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type Commands struct {
	cl  api.ItisaDBClient
	sdk *itisadb.Client
}

func New(cl api.ItisaDBClient, sdk *itisadb.Client) *Commands {
	return &Commands{
		cl:  cl,
		sdk: sdk,
	}
}

type Action string

var ErrWrongInput = errors.New("wrong input")
var ErrUnknownCMD = errors.New("unknown cmd")
var ErrEmpty = errors.New("the value does not exist")
var ErrUnknownServer = errors.New("the value does not exist")

const (
	get      = "get"
	newEl    = "new"
	object   = "object"
	marshalo = "marshalo"
	attach   = "attach"
	seto     = "seto"
	geto     = "geto"
	delo     = "delo"
	del      = "del"
	deleteEl = "delete"
)

const (
	level = "level"
	on    = "on"
)

const (
	_ = iota * -1
	setToAll
)

func (c *Commands) Do(ctx context.Context, act string, args ...string) (string, error) {
	switch strings.ToLower(act) {
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
	case Set:
		cmd, err := ParseSet(args)
		if err != nil {
			return "", err
		}

		return c.set(ctx, cmd)
	case newEl:
		if len(args) < 1 {
			return "", ErrWrongInput
		}

		switch strings.ToLower(args[0]) {
		case object:
			if len(args) < 2 || len(args) > 4 {
				return "", ErrWrongInput
			}
			name := args[1]

			var server int32
			var lvl int8

			checkIsLevel := func(probablyLevel string) bool {
				if probablyLevel == "S" || probablyLevel == "s" {
					return true
				}
				if probablyLevel == "R" || probablyLevel == "r" {
					return true
				}
				return false
			}

			for i := 2; i < len(args); i++ {
				if checkIsLevel(args[i]) {
					if args[i] == "S" {
						lvl = secretLevel
					} else if args[i] == "R" {
						lvl = restrictedLevel
					} else {
						return "", fmt.Errorf("unknown level: %s", args[i])
					}
				} else {
					serverStr := args[i]
					serverInt, err := strconv.Atoi(serverStr)
					if err != nil {
						return "", err
					}
					server = int32(serverInt)
				}
			}

			switch r := c.sdk.Object(ctx, name, itisadb.ObjectOptions{
				Server: &server,
				Level:  itisadb.Level(lvl),
			}); r.Switch() {
			case gost.IsOk:
				return fmt.Sprintf("status: ok, object %s created", name), nil
			case gost.IsErr:
				return "", r.Error()
			}

			return "", ErrWrongInput
		default:
			return "", ErrWrongInput
		}
	case seto:
		sc, err := ParseSet(args[1:])
		if err != nil {
			return "", err
		}
		return c.setObject(ctx, args[0], sc.key, sc.value)
	case geto:
		switch len(args) {
		case 1:
			return c.showObject(ctx, args[0])
		case 2:
			return c.getFromObject(ctx, args[0], args[1])
		default:
			return "", ErrWrongInput
		}
	case marshalo:
		if len(args) < 1 {
			return "", ErrWrongInput
		}
		name := args[0]
		return c.showObject(ctx, name)
	case del:
		if len(args) < 1 {
			return "", ErrWrongInput
		}
		name := args[0]
		return c.del(ctx, name)
	case delo:
		if len(args) < 2 {
			return "", ErrWrongInput
		}
		object := args[0]
		key := args[1]
		return c.delo(ctx, object, key)
	case deleteEl:
		if len(args) < 2 {
			return "", ErrWrongInput
		}

		switch strings.ToLower(args[0]) {
		case object:
			object := args[1]
			return c.deleteObject(ctx, object)
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

func (c *Commands) showObject(ctx context.Context, name string) (string, error) {
	m, err := c.cl.ObjectToJSON(ctx, &api.ObjectToJSONRequest{Name: name})
	if err != nil {
		return "", err
	}

	return m.Object, nil
}

func (c *Commands) object(ctx context.Context, act, name, key, value string) (string, error) {
	switch act {
	case seto:
		return c.setObject(ctx, name, key, value)
	case geto:
		return c.getFromObject(ctx, name, key)
	case marshalo:
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

func (c *Commands) get(ctx context.Context, key string, server int32) (string, error) {
	switch r := c.sdk.GetOne(ctx, key, itisadb.GetOptions{Server: &server}); r.Switch() {
	case gost.IsOk:
		return r.Unwrap(), nil
	case gost.IsErr:
		return "", fmt.Errorf(r.Error().Error())
	}

	return "", nil
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

func toServerNumber(server int32) *int32 {
	if server != 0 {
		return &server
	}

	return nil
}

func (c *Commands) set(ctx context.Context, cmd Command) (string, error) {
	args := cmd.Args()

	if len(args) < 2 {
		return "", ErrWrongInput
	}

	r := c.sdk.SetOne(ctx, args[0], args[1], itisadb.SetOptions{
		Server:   toServerNumber(cmd.Server()),
		ReadOnly: cmd.Mode() == readOnlySetMode,
		Level:    itisadb.Level(cmd.Level()),
		Unique:   cmd.Mode() == uniqueSetMode,
	})

	if err := r.Error(); err != nil {
		return "", fmt.Errorf(err.Error())
	}

	resp := ""

	switch savedTo := r.Unwrap(); savedTo {
	case 0:
		resp = fmt.Sprintf("status: ok, saved on server #%d", savedTo)
	case setToAll:
		resp = fmt.Sprintf("status: ok, saved on all servers")
	default:
		resp = fmt.Sprintf("status: ok, saved on server #%d", savedTo)
	}

	return resp, nil
}

func (c *Commands) getFromObject(ctx context.Context, name string, key string) (string, error) {
	r, err := c.cl.GetFromObject(ctx, &api.GetFromObjectRequest{
		Object: name,
		Key:    key,
	})

	if err != nil {
		return "", err
	}

	return r.Value, nil
}

func (c *Commands) del(ctx context.Context, name string) (string, error) {
	_, err := c.cl.Delete(ctx, &api.DeleteRequest{Key: name})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, deleted key %s", name), nil
}

func (c *Commands) delo(ctx context.Context, object, key string) (string, error) {
	_, err := c.cl.DeleteAttr(ctx, &api.DeleteAttrRequest{Key: key, Object: object})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, deleted object %s attribute %s", object, key), nil
}

func (c *Commands) deleteObject(ctx context.Context, object string) (string, error) {
	_, err := c.cl.DeleteObject(ctx, &api.DeleteObjectRequest{Object: object})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("status: ok, deleted object %s", object), nil
}

type Command interface {
	Action() string
	Args() []string
	Server() int32
	Mode() uint8
	Level() uint8
}

func ParseCommand(ctx context.Context, text string) (Command, error) {
	split := strings.Split(text, " ")
	if len(split) < 1 {
		return nil, fmt.Errorf("no command detected")
	}

	switch cmd := strings.ToLower(split[0]); cmd {
	case Set:
		return ParseSet(split[1:])
	}

	return nil, fmt.Errorf("unknown command")
}
