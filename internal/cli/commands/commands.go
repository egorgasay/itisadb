package commands

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	"strconv"
	"strings"
)

type Commands struct {
	sdk *itisadb.Client
}

func New(sdk *itisadb.Client) *Commands {
	return &Commands{
		sdk: sdk,
	}
}

const (
	UnknownCode = iota
	InvalidCode
)

const (
	UnknownExtCode = iota
	InputExtCode
	CmdExtCode
)

type Action string

var (
	ErrWrongInput = gost.NewError(InvalidCode, InputExtCode, "wrong input")
	ErrUnknownCMD = gost.NewError(UnknownCode, CmdExtCode, "unknown cmd")
)

const (
	_get      = "get"
	_newEl    = "new"
	_object   = "object"
	_marshalo = "marshalo"
	_attach   = "attach"
	_seto     = "seto"
	_geto     = "geto"
	_delo     = "delo"
	_del      = "del"
	_deleteEl = "delete"
)

const (
	_ = iota * -1
	setToAll
)

func (c *Commands) Do(ctx context.Context, act string, args ...string) (res gost.Result[string]) {
	switch strings.ToLower(act) {
	case _get:
		return c.get(ctx, args)
	case Set:
		cmd, err := ParseSet(args)
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, err.Error())
		}

		r := c.set(ctx, cmd)
		if r.IsErr() {
			return res.ErrNew(InvalidCode, InputExtCode, r.Error().Error())
		}

		switch savedTo := r.Unwrap(); savedTo {
		case setToAll:
			res.Ok(fmt.Sprintf("status: ok, saved on all servers"))
		default:
			res.Ok(fmt.Sprintf("status: ok, saved on server #%d", savedTo))
		}
	case _newEl:
		if len(args) < 1 {
			return res.Err(ErrWrongInput)
		}

		switch strings.ToLower(args[0]) {
		case _object:
			r := c.newObject(ctx, args[1:])
			if r.IsOk() {
				return res.Ok(fmt.Sprintf("object %s created", r.Unwrap().Name()))
			}
			return res.ErrNew(InvalidCode, InputExtCode, r.Error().Error())
		default:
			return res.Err(ErrWrongInput)
		}
	case _seto:
		sc, err := ParseSet(args[1:])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, err.Error())
		}

		if r := c.setObject(ctx, args[0], sc); r.IsOk() {
			return res.Ok(fmt.Sprintf("status: ok, saved in object %s, on server #%d", args[0], r.Unwrap()))
		} else {
			return res.Err(r.Error())
		}
	case _geto:
		r := c.geto(ctx, args)
		return r
	case _marshalo:
		if len(args) < 1 {
			return res.Err(ErrWrongInput)
		}
		name := args[0]
		r := c.showObject(ctx, name)
		return r
	case _del:
		if len(args) < 1 {
			return res.Err(ErrWrongInput)
		}

		name := args[0]

		var opts itisadb.DeleteOptions
		if len(args) >= 2 {
			serverInt, err := strconv.Atoi(args[1])
			if err != nil {
				return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[1]))
			}

			serverI32 := int32(serverInt)
			opts.Server = &serverI32
		}

		if r := c.sdk.DelOne(ctx, name, opts); r.IsOk() {
			return res.Ok(fmt.Sprintf("status: ok, deleted %s", name))
		} else {
			return res.Err(r.Error())
		}
	case _delo:
		if len(args) < 2 {
			return res.Err(ErrWrongInput)
		}

		var opts itisadb.ObjectOptions
		if len(args) >= 2 {
			serverInt, err := strconv.Atoi(args[1])
			if err != nil {
				return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[1]))
			}

			serverI32 := int32(serverInt)
			opts.Server = &serverI32
		}

		obj := args[0]
		key := args[1]

		r := c.sdk.Object(ctx, obj, opts)
		if r.IsErr() {
			return res.Err(r.Error())
		}

		if r := r.Unwrap().DeleteKey(ctx, key); r.IsErr() {
			return res.Err(r.Error())
		}

		return res.Ok(fmt.Sprintf("status: ok, deleted %s from %s", key, obj))
	case _deleteEl:
		if len(args) < 2 {
			return res.Err(ErrWrongInput)
		}

		switch strings.ToLower(args[0]) {
		case _object:
			var aargs []string

			object := args[1]
			if len(args) >= 3 {
				aargs = args[2:]
			} else {
				aargs = make([]string, 0)
			}

			r := c.deleteObject(ctx, object, aargs)
			if r.IsErr() {
				return res.Err(r.Error())
			}
			return res.Ok(fmt.Sprintf("status: ok, deleted %s", object))
		}
	case _attach:
		if len(args) < 2 {
			return res.Err(ErrWrongInput)
		}

		dst := args[0]
		src := args[1]

		if r := c.attach(ctx, dst, src); r.IsErr() {
			return res.Err(r.Error())
		}

		return res.Ok(fmt.Sprintf("status: ok, attached %s to %s", src, dst))
	}

	return res.Err(ErrUnknownCMD)
}

func (c *Commands) newObject(ctx context.Context, args []string) (res gost.Result[*itisadb.Object]) {
	if len(args) < 1 || len(args) > 3 {
		return res.Err(ErrWrongInput)
	}

	name := args[1]

	var (
		server int32
		lvl    int8
	)

	for i := 2; i < len(args); i++ {
		if args[i] == "S" || args[i] == "s" {
			lvl = secretLevel
		} else if args[i] == "R" || args[i] == "r" {
			lvl = restrictedLevel
		} else {
			serverStr := args[i]

			serverInt, err := strconv.Atoi(serverStr)
			if err != nil {
				return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", serverStr))
			}

			server = int32(serverInt)
		}
	}

	switch r := c.sdk.Object(ctx, name, itisadb.ObjectOptions{
		Server: &server,
		Level:  itisadb.Level(lvl),
	}); r.Switch() {
	case gost.IsOk:
		return res.Ok(r.Unwrap())
	case gost.IsErr:
		return res.Err(r.Error())
	}

	return res.Err(ErrWrongInput)
}

func (c *Commands) geto(ctx context.Context, args []string) (res gost.Result[string]) {
	switch len(args) {
	case 1:
		return c.showObject(ctx, args[0])
	case 2:
		return c.getFromObject(ctx, args[0], args[1])
	default:
		return res.Err(ErrWrongInput)
	}
}

func (c *Commands) showObject(ctx context.Context, name string) (res gost.Result[string]) {
	if r := c.sdk.Object(ctx, name); r.IsOk() {
		return r.Unwrap().JSON(ctx)
	} else {
		return res.Err(r.Error())
	}
}

func (c *Commands) setObject(ctx context.Context, object string, cmd SetCommand) (res gost.Result[int32]) {
	if r := c.sdk.Object(ctx, object); r.IsOk() {
		return r.Unwrap().Set(ctx, cmd.key, cmd.value)
	} else {
		return res.Err(r.Error())
	}
}

func (c *Commands) get(ctx context.Context, args []string) (res gost.Result[string]) {
	if len(args) < 1 {
		return res.Err(ErrWrongInput)
	}

	var server *int32
	if len(args) != 1 {
		num, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[len(args)-1]))
		}
		serverInt32 := int32(num)

		server = &serverInt32
	}

	return c.sdk.GetOne(ctx, args[0], itisadb.GetOptions{Server: server})
}

func (c *Commands) attach(ctx context.Context, dst string, src string) (res gost.Result[gost.Nothing]) {
	r := c.sdk.Object(ctx, dst)
	if r.IsOk() {
		return r.Unwrap().Attach(ctx, src)
	} else {
		return res.Err(r.Error())
	}
}

func toServerNumber(server int32) *int32 {
	if server != 0 {
		return &server
	}

	return nil
}

func (c *Commands) set(ctx context.Context, cmd Command) (res gost.Result[int32]) {
	args := cmd.Args()

	if len(args) < 2 {
		return res.Err(ErrWrongInput)
	}

	return c.sdk.SetOne(ctx, args[0], args[1], itisadb.SetOptions{
		Server:   toServerNumber(cmd.Server()),
		ReadOnly: cmd.Mode() == readOnlySetMode,
		Level:    itisadb.Level(cmd.Level()),
		Unique:   cmd.Mode() == uniqueSetMode,
	})
}

func (c *Commands) getFromObject(ctx context.Context, name string, key string, server ...string) (res gost.Result[string]) {
	opts := itisadb.GetFromObjectOptions{}

	if len(server) == 1 {
		servInt, err := strconv.Atoi(server[0])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("invalid server number %s", server[0]))
		}

		serv := int32(servInt)
		opts.Server = &serv
	}

	if r := c.sdk.Object(ctx, name); r.IsOk() {
		return r.Unwrap().Get(ctx, key, opts)
	} else {
		return res.Err(r.Error())
	}
}

func (c *Commands) deleteObject(ctx context.Context, object string, args []string) (res gost.Result[gost.Nothing]) {
	if len(args) < 1 {
		return res.Err(ErrWrongInput)
	}

	var server *int32
	{
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[0]))
		}
		serverInt32 := int32(num)
		server = &serverInt32
	}

	r := c.sdk.Object(ctx, object, itisadb.ObjectOptions{Server: server})
	if r.IsErr() {
		return res.Err(r.Error())
	}

	return r.Unwrap().DeleteObject(ctx)
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
