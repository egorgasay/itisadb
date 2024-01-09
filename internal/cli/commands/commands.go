package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
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
	_marshalo = "marshalo"
	_attach   = "attach"
	_seto     = "seto"
	_geto     = "geto"
	_delo     = "delo"
	_del      = "del"
	_deleteEl = "delete"
	_change   = "change"

	_object       = "object"
	_userLevel    = "user.level"
	_userPassword = "user.password"
	_user         = "user"
)

const (
	_ = iota * -1
	setToAll
)

func (c *Commands) Do(ctx context.Context, act string, args ...string) (res gost.Result[string]) {
	switch strings.ToLower(act) {
	case _get:
		r := c.get(ctx, args)
		if r.IsErr() {
			return res.Err(r.Error())
		}

		b, err := json.MarshalIndent(r.Unwrap(), "&ensp;", "&ensp;") // TODO: redo
		if err != nil {
			return res.ErrNew(0, 0, fmt.Sprintf("cannot marshal: %s", err.Error()))
		}

		return res.Ok(string(b))
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
			return res.Ok(fmt.Sprintf("status: ok, saved on all balancer"))
		default:
			return res.Ok(fmt.Sprintf("status: ok, saved on server #%d", savedTo))
		}
	case _newEl:
		if len(args) < 1 {
			return res.Err(ErrWrongInput)
		}

		switch strings.ToLower(args[0]) {
		case _user: // NEW USER <name> <password> <optional level>
			if len(args) < 3 {
				return res.Err(ErrWrongInput)
			}

			opts := itisadb.NewUserOptions{}

			if len(args) >= 4 {
				levelRes := levelFromStr(args[len(args)-1])
				if res.IsErr() {
					return res.Err(levelRes.Error())
				}

				opts.Level = levelRes.Unwrap()
			}

			r := c.sdk.NewUser(ctx, args[1], args[2], opts)
			if r.IsErr() {
				return res.Err(r.Error())
			}

			return res.Ok(fmt.Sprintf("user %s with level %s created", args[1], opts.Level.String()))
		case _object:
			r := c.newObject(ctx, args[1:])
			if r.IsOk() {
				obj := r.Unwrap()
				return res.Ok(fmt.Sprintf("object %s created on server #%d", obj.Name(), obj.Server()))
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
			server, err := strconv.Atoi(args[1])
			if err != nil {
				return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[1]))
			}

			opts.Server = int32(server)
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
			server, err := strconv.Atoi(args[1])
			if err != nil {
				return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[1]))
			}

			opts.Server = int32(server)
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
		case _user:
			user := args[1]

			r := c.sdk.DeleteUser(ctx, user)
			if r.IsErr() {
				return res.Err(r.Error())
			}

			return res.Ok(fmt.Sprintf("status: ok, deleted %s", user))
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
	case _change:
		if len(args) < 2 {
			return res.Err(ErrWrongInput)
		}

		switch strings.ToLower(args[0]) {
		case _userLevel: // CHANGE user.level <user> <level>
			args = args[1:]

			if len(args) < 2 {
				return res.Err(ErrWrongInput)
			}

			user := args[0]
			levelRes := levelFromStr(args[1])

			if levelRes.IsErr() {
				return res.Err(levelRes.Error())
			}

			level := levelRes.Unwrap()

			resLevel := c.sdk.ChangeLevel(ctx, user, level)
			if resLevel.IsErr() {
				return res.Err(resLevel.Error())
			}

			return res.Ok(fmt.Sprintf("status: ok, changed %s level to %s", user, level))
		case _userPassword: // CHANGE user.password <user> <password>
			args = args[1:]

			if len(args) < 2 {
				return res.Err(ErrWrongInput)
			}

			user := args[0]
			password := args[1]

			resPassword := c.sdk.ChangePassword(ctx, user, password)
			if resPassword.IsErr() {
				return res.Err(resPassword.Error())
			}

			return res.Ok(fmt.Sprintf("status: ok, changed %s password", user))
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

func levelFromStr(lvl string) (res gost.Result[itisadb.Level]) {
	lvl = strings.ToLower(lvl)
	if lvl == "s" {
		return res.Ok(secretLevel)
	} else if lvl == "r" {
		return res.Ok(restrictedLevel)
	}

	return res.Err(ErrWrongInput.WrapfNotNilMsg("wrong level: %s", lvl))
}

func (c *Commands) newObject(ctx context.Context, args []string) (res gost.Result[*itisadb.Object]) {
	if len(args) < 1 || len(args) > 3 {
		return res.Err(ErrWrongInput)
	}

	name := args[0]

	var (
		server int32
		lvl    int8
	)

	for i := 1; i < len(args); i++ {
		word := strings.ToLower(args[i])
		if word == "s" {
			lvl = secretLevel
		} else if word == "r" {
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
		Server: server,
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

func (c *Commands) get(ctx context.Context, args []string) (res gost.Result[itisadb.Value]) {
	if len(args) < 1 {
		return res.Err(ErrWrongInput)
	}

	var server int32
	if len(args) != 1 {
		num, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[len(args)-1]))
		}

		server = int32(num)
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
		Server:   cmd.Server(),
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

		opts.Server = int32(servInt)
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

	var server int32
	{
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return res.ErrNew(InvalidCode, InputExtCode, fmt.Sprintf("wrong server number: %s", args[0]))
		}
		server = int32(num)
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
