package balancer

import (
	"context"
	"fmt"
	"strings"

	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

func (c *Balancer) Object(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.ObjectOptions) (s int32, err error) {
	return s, gost.WithContextPool(ctx, func() error {
		s, err = c.object(ctx, claims, name, opts)
		return err
	}, c.pool)
}

func (c *Balancer) addObjectServer(object string, server int32) {
	objects := strings.Split(object, constants.ObjectSeparator)
	if len(objects) > 0 {
		object = objects[0]
	}

	defer c.objectServers.WRelease()
	(*c.objectServers.WBorrow().Ref())[object] = server
}

func (c *Balancer) getObjectServer(object string) (opt gost.Option[int32]) {
	objects := strings.Split(object, constants.ObjectSeparator)
	if len(objects) > 0 {
		object = objects[0]
	}

	defer c.objectServers.Release()

	server, ok := c.objectServers.RBorrow().Read()[object]
	if !ok {
		return opt.None()
	}

	return opt.Some(server)
}

func (c *Balancer) delObjectServer(object string) {
	defer c.objectServers.WRelease()
	delete(c.objectServers.WBorrow().Read(), object)
}

func (c *Balancer) addKeyServer(key string, server int32) {
	defer c.keyServers.WRelease()
	(*c.keyServers.WBorrow().Ref())[key] = server
}

func (c *Balancer) getKeyServer(key string) (opt gost.Option[int32]) {
	defer c.keyServers.Release()

	server, ok := c.keyServers.RBorrow().Read()[key]
	if !ok {
		return opt.None()
	}

	return opt.Some(server)
}

func (c *Balancer) delKeyServer(key string) {
	defer c.keyServers.WRelease()
	delete(*c.keyServers.WBorrow().Ref(), key)
}

func (c *Balancer) isObject(ctx context.Context, claims gost.Option[models.UserClaims], object string) (res gost.Result[gost.Option[int32]]) {
	var v gost.Option[int32]

	err := c.servers.Iter(func(server domains.Server) error {
		if r := server.IsObject(ctx, claims, object, models.IsObjectOptions{}); r.IsErr() {
			return r.Error()
		} else {
			if r.Unwrap() {
				v = v.Some(server.Number())
			}
		}

		return nil
	})

	if err != nil {
		return res.ErrNewUnknown(err.Error())
	}

	return res.Ok(v)
}

func (c *Balancer) object(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.ObjectOptions) (int32, error) {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return 0, r.Error()
	}

	serv := r.Unwrap()

	if r := serv.NewObject(ctx, claims, object, opts); r.IsErr() {
		return 0, fmt.Errorf("can't create object: %w", r.Error())
	}

	c.addObjectServer(object, serv.Number())

	return serv.Number(), nil
}

func (c *Balancer) GetFromObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key string, opts models.GetFromObjectOptions) (v string, err error) {
	return v, gost.WithContextPool(ctx, func() error {
		v, err = c.getFromObject(ctx, claims, object, key, opts)
		return err
	}, c.pool)
}

func (c *Balancer) findServerForObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, server int32) (res gost.Result[domains.Server]) {
	objects := strings.Split(object, constants.ObjectSeparator)
	var resolvedServer int32

	serverOpt := c.getObjectServer(objects[0])
	if serverOpt.IsNone() {
		r := c.isObject(ctx, claims, object)
		if r.IsErr() {
			return res.Err(r.Error())
		}

		if opt := r.Unwrap(); opt.IsSome() {
			resolvedServer = opt.Unwrap()
		}
	} else {
		resolvedServer = serverOpt.Unwrap()
	}

	isRequestedServerAuto := server == constants.AutoServerNumber
	isResolvedServerNone := resolvedServer == 0

	if !isRequestedServerAuto && !isResolvedServerNone && server != resolvedServer {
		return res.Err(constants.ErrAlreadyExists.ExtendMsg(fmt.Sprintf("can't get inner object[%d] from different[%d] server", server, serverOpt)))
	}

	s, ok := c.servers.GetServer(resolvedServer)
	if !ok {
		return res.Err(constants.ErrServerNotFound)
	}

	return res.Ok(s)
}

func (c *Balancer) getFromObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key string, opts models.GetFromObjectOptions) (string, error) {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return "", r.Error()
	}

	if r := r.Unwrap().GetFromObject(ctx, claims, object, key, opts); r.IsErr() {
		return "", r.Error()
	} else {
		return r.Unwrap(), nil
	}
}

func (c *Balancer) SetToObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key, val string, opts models.SetToObjectOptions) (s int32, err error) {
	return s, gost.WithContextPool(ctx, func() error {
		s, err = c.setToObject(ctx, claims, object, key, val, opts)
		return err
	}, c.pool)
}

func (c *Balancer) setToObject(ctx context.Context, claims gost.Option[models.UserClaims], object, key, val string, opts models.SetToObjectOptions) (int32, error) {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return 0, r.Error()
	}

	cl := r.Unwrap()

	rSet := cl.SetToObject(ctx, claims, object, key, val, opts)
	if rSet.IsErr() {
		return 0, fmt.Errorf("can't set object: %w", rSet.Error())
	}

	c.addObjectServer(object, cl.Number())

	return cl.Number(), nil
}

func (c *Balancer) ObjectToJSON(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.ObjectToJSONOptions) (string, error) {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return "", r.Error()
	}

	cl := r.Unwrap()

	resObj := cl.ObjectToJSON(ctx, claims, object, opts)
	if resObj.IsErr() {
		return "", resObj.Error()
	}

	return resObj.Unwrap(), nil
}

func (c *Balancer) IsObject(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.IsObjectOptions) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	// TODO:
	return c.getObjectServer(name).IsSome(), nil
}

func (c *Balancer) Size(ctx context.Context, claims gost.Option[models.UserClaims], name string, opts models.SizeOptions) (size uint64, err error) {
	return size, gost.WithContextPool(ctx, func() error {
		size, err = c.size(ctx, claims, name, opts)
		return err
	}, c.pool)
}

func (c *Balancer) size(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.SizeOptions) (uint64, error) {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return 0, r.Error()
	}

	res := r.Unwrap().ObjectSize(ctx, claims, object, opts)
	return res.UnwrapOrDefault(), res.ErrorStd()
}

func (c *Balancer) DeleteObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.DeleteObjectOptions) error {
	return gost.WithContextPool(ctx, func() error {
		return c.deleteObject(ctx, claims, object, opts)
	}, c.pool)
}

func (c *Balancer) deleteObject(ctx context.Context, claims gost.Option[models.UserClaims], object string, opts models.DeleteObjectOptions) error {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return r.Error()
	}

	if r := r.Unwrap().DeleteObject(ctx, claims, object, opts); r.IsErr() {
		return fmt.Errorf("can't delete object: %w", r.Error())
	}

	c.delObjectServer(object)

	return nil
}

func (c *Balancer) AttachToObject(ctx context.Context, claims gost.Option[models.UserClaims], dst, src string, opts models.AttachToObjectOptions) error {
	return gost.WithContextPool(ctx, func() error {
		return c.attachToObject(ctx, claims, dst, src, opts)
	}, c.pool)
}

func (c *Balancer) attachToObject(ctx context.Context, claims gost.Option[models.UserClaims], dst, src string, opts models.AttachToObjectOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(dst)
		if res.IsNone() {
			return fmt.Errorf("can't get dst object server: %w", constants.ErrObjectNotFound)
		}

		server1 := res.Unwrap()

		res = c.getObjectServer(src)
		if res.IsNone() {
			return fmt.Errorf("can't get src object server: %w", constants.ErrObjectNotFound)
		}

		server2 := res.Unwrap()

		if server1 != server2 {
			return fmt.Errorf("can't attach to objects from different servers: %w", constants.ErrForbidden)
		}

		opts.Server = server1
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	if r := cl.AttachToObject(ctx, claims, dst, src, opts); r.IsErr() {
		return fmt.Errorf("can't attach to object: %w", r.Error())
	}

	return nil
}

func (c *Balancer) DeleteAttr(ctx context.Context, claims gost.Option[models.UserClaims], key string, object string, opts models.DeleteAttrOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return gost.WithContextPool(ctx, func() error {
		return c.deleteAttr(ctx, claims, key, object, opts)
	}, c.pool)
}

func (c *Balancer) deleteAttr(ctx context.Context, claims gost.Option[models.UserClaims], key, object string, opts models.DeleteAttrOptions) error {
	r := c.findServerForObject(ctx, claims, object, opts.Server)
	if r.IsErr() {
		return r.Error()
	}

	if r := r.Unwrap().ObjectDeleteKey(ctx, claims, key, object, opts); r.IsErr() {
		return fmt.Errorf("can't delete attr: %w", r.Error())
	}

	return nil
}
