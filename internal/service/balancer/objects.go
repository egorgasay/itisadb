package balancer

import (
	"context"
	"fmt"
	"github.com/egorgasay/gost"
	"itisadb/internal/constants"
	"itisadb/internal/domains"
	"itisadb/internal/models"
)

func (c *Balancer) Object(ctx context.Context, userID int, name string, opts models.ObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.object(ctx, userID, name, opts)
		return err
	})
}

func (c *Balancer) addObjectServer(object string, server int32) {
	defer c.objectServers.WRelease()
	(*c.objectServers.WBorrow().Ref())[object] = server
}

func (c *Balancer) getObjectServer(object string) (opt gost.Option[int32]) {
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
	delete(c.keyServers.WBorrow().Read(), key)
}

func (c *Balancer) object(ctx context.Context, userID int, name string, opts models.ObjectOptions) (int32, error) {
	if opts.Server == constants.AutoServerNumber {
		if res := c.getObjectServer(name); res.IsSome() {
			return res.Unwrap(), nil
		}
	}

	var serv domains.Server
	var ok bool

	serv, ok = c.servers.GetServer(opts.Server)
	if !ok {
		return 0, constants.ErrServerNotFound
	}

	r := serv.NewObject(ctx, userID, name, opts)
	if r.IsErr() {
		return 0, fmt.Errorf("can't create object: %w", r.Error())
	}

	c.addObjectServer(name, opts.Server)

	return serv.Number(), nil
}

func (c *Balancer) GetFromObject(ctx context.Context, userID int, object, key string, opts models.GetFromObjectOptions) (v string, err error) {
	return v, c.withContext(ctx, func() error {
		v, err = c.getFromObject(ctx, userID, object, key, opts)
		return err
	})
}

func (c *Balancer) getFromObject(ctx context.Context, userID int, object, key string, opts models.GetFromObjectOptions) (string, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return "", constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	if r := cl.GetFromObject(ctx, userID, object, key, opts); r.IsErr() {
		return "", r.Error()
	} else {
		return r.Unwrap(), nil
	}
}

func (c *Balancer) SetToObject(ctx context.Context, userID int, object, key, val string, opts models.SetToObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.setToObject(ctx, userID, object, key, val, opts)
		return err
	})
}

func (c *Balancer) setToObject(ctx context.Context, userID int, object, key, val string, opts models.SetToObjectOptions) (int32, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return 0, constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	r := cl.SetToObject(ctx, userID, object, key, val, opts)
	if r.IsErr() {
		return 0, fmt.Errorf("can't set object: %w", r.Error())
	}

	return cl.Number(), nil
}

func (c *Balancer) ObjectToJSON(ctx context.Context, userID int, object string, opts models.ObjectToJSONOptions) (string, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return "", constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	resObj := cl.ObjectToJSON(ctx, userID, object, opts)
	if resObj.IsErr() {
		return "", resObj.Error()
	}

	return resObj.Unwrap(), nil
}

func (c *Balancer) IsObject(ctx context.Context, userID int, name string, opts models.IsObjectOptions) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return c.getObjectServer(name).IsSome(), nil
}

func (c *Balancer) Size(ctx context.Context, userID int, name string, opts models.SizeOptions) (size uint64, err error) {
	return size, c.withContext(ctx, func() error {
		size, err = c.size(ctx, userID, name, opts)
		return err
	})
}

func (c *Balancer) size(ctx context.Context, userID int, object string, opts models.SizeOptions) (uint64, error) {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return 0, constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	res := cl.ObjectSize(ctx, userID, object, opts)
	return res.UnwrapOrDefault(), res.Error()
}

func (c *Balancer) DeleteObject(ctx context.Context, userID int, object string, opts models.DeleteObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, userID, object, opts)
	})
}

func (c *Balancer) deleteObject(ctx context.Context, userID int, object string, opts models.DeleteObjectOptions) error {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	if r := cl.DeleteObject(ctx, userID, object, opts); r.IsErr() {
		return fmt.Errorf("can't delete object: %w", r.Error())
	}

	c.delObjectServer(object)

	return nil
}

func (c *Balancer) AttachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.attachToObject(ctx, userID, dst, src, opts)
	})
}

func (c *Balancer) attachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) error {
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

	return cl.AttachToObject(ctx, userID, dst, src, opts).Error().WrapNotNilMsg("can't attach to object")
}

func (c *Balancer) DeleteAttr(ctx context.Context, userID int, key string, object string, opts models.DeleteAttrOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, userID, key, object, opts)
	})
}

func (c *Balancer) deleteAttr(ctx context.Context, userID int, key, object string, opts models.DeleteAttrOptions) error {
	if opts.Server == constants.AutoServerNumber {
		res := c.getObjectServer(object)
		if res.IsNone() {
			return constants.ErrObjectNotFound
		}

		opts.Server = res.Unwrap()
	}

	cl, ok := c.servers.GetServer(opts.Server)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	if r := cl.ObjectDeleteKey(ctx, userID, key, object, opts); r.IsErr() {
		return fmt.Errorf("can't delete attr: %w", r.Error())
	}

	return nil
}
