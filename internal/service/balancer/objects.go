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
	if !c.hasPermission(userID, opts.Level) {
		return 0, constants.ErrForbidden
	}

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

func (c *Balancer) object(ctx context.Context, userID int, name string, opts models.ObjectOptions) (int32, error) {
	if opts.Server == constants.AutoServerNumber {
		if res := c.getObjectServer(name); res.IsSome() {
			return res.Unwrap(), nil
		}
	}

	var serv domains.Server
	var ok bool

	serv, ok = c.servers.GetServerByID(opts.Server)
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
	cl, ok := c.servers.GetServerByID(opts.Server)
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
	res := c.getObjectServer(object)
	if res.IsNone() {
		return 0, constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(res.Unwrap())
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	r := cl.SetToObject(ctx, userID, object, key, val, opts)
	if r.IsErr() {
		return 0, fmt.Errorf("can't set object: %w", r.Error())
	}

	return cl.Number(), nil
}

func (c *Balancer) ObjectToJSON(ctx context.Context, userID int, name string, opts models.ObjectToJSONOptions) (string, error) {
	cl, ok := c.servers.GetServerByID(opts.Server)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	res := cl.ObjectToJSON(ctx, name, opts)
	if res.IsErr() {
		return "", res.Error()
	}

	return res.Unwrap(), nil
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

func (c *Balancer) size(ctx context.Context, userID int, name string, opts models.SizeOptions) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	info, err := c.storage.GetObjectInfo(name)
	if err != nil {
		return 0, fmt.Errorf("can't get object info: %w", err)
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return 0, constants.ErrServerNotFound
		}

		//res, err := cl.Size(ctx, name, opts)
		//if err != nil {
		//	return 0, err 	TODO:
		//}
		//
		//return res.Size, nil

		return 0, nil // TODO:
	}

	if info.Server != _mainStorage {
		return 0, constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return 0, constants.ErrForbidden
	}

	size, err := c.storage.Size(name)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func (c *Balancer) DeleteObject(ctx context.Context, userID int, name string, opts models.DeleteObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, userID, name, opts)
	})
}

func (c *Balancer) deleteObject(ctx context.Context, userID int, name string, opts models.DeleteObjectOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	info, err := c.storage.GetObjectInfo(name)
	if err != nil {
		return fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return constants.ErrServerNotFound
		}

		//err = cl.DelOne(ctx, name, opts.ToSDK()).Error()
		//if err != nil {
		//	return err
		//}
		//
		//c.storage.DeleteObjectInfo(name)
		//
		//return nil

		return nil
	}

	if info.Server != _mainStorage {
		return constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return constants.ErrForbidden
	}

	err = c.storage.DeleteObject(name)
	if err != nil {
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteDeleteObject(name)
	}

	c.storage.DeleteObjectInfo(name)

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

	info, err := c.storage.GetObjectInfo(src)
	if err != nil {
		return fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return constants.ErrServerNotFound
		}

		//err = cl.AttachToObject(ctx, dst, src, opts)
		//if err != nil {
		//	return err
		//} TODO:
		//
		//return nil

		return nil // TODO:
	}

	if info.Server != _mainStorage {
		return constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return constants.ErrForbidden
	}

	err = c.storage.AttachToObject(dst, src)
	if err != nil {
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteAttach(dst, src)
	}

	return nil
}

func (c *Balancer) DeleteAttr(ctx context.Context, userID int, key string, object string, opts models.DeleteAttrOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, userID, key, object, opts)
	})
}

func (c *Balancer) deleteAttr(ctx context.Context, userID int, key, object string, opts models.DeleteAttrOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return constants.ErrServerNotFound
		}

		//err = cl.DeleteAttr(ctx, key, object, opts)
		//if err != nil {
		//	return err
		//} TODO:

		return nil
	}

	if info.Server != _mainStorage {
		return constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return constants.ErrForbidden
	}

	err = c.storage.DeleteAttr(object, key)
	if err != nil {
		return err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteDeleteAttr(object, key)
	}

	return nil
}
