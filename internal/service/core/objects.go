package core

import (
	"context"
	"errors"
	"fmt"
	"itisadb/internal/constants"
	"itisadb/internal/models"
	"itisadb/internal/service/servers"
	"itisadb/pkg"
)

func (c *Core) earlyObjectNotFound(requested *int32, actual int32) bool {
	if requested != nil {
		return *requested != actual
	}

	return false
}

func (c *Core) Object(ctx context.Context, userID int, name string, opts models.ObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.object(ctx, userID, name, opts)
		return err
	})
}

func (c *Core) object(ctx context.Context, userID int, name string, opts models.ObjectOptions) (int32, error) {
	info, err := c.storage.GetObjectInfo(name)
	errNotFound := errors.Is(err, constants.ErrObjectNotFound)

	if err != nil && !errNotFound {
		return 0, fmt.Errorf("can't get object info: %w", err)
	}

	if !c.useMainStorage(opts.Server) {
		var serv *servers.Server
		var ok bool

		if errNotFound {
			serv, ok = c.servers.GetServer()
		} else {
			serv, ok = c.servers.GetServerByID(info.Server)
		}

		if !ok {
			return 0, constants.ErrServerNotFound
		}

		err := serv.NewObject(ctx, name, opts)
		if err != nil {
			return 0, fmt.Errorf("can't create object: %w", err)
		}

		num := serv.GetNumber()

		info := models.ObjectInfo{
			Server: num,
			Level:  opts.Level,
		}

		if c.cfg.TransactionLogger.On {
			c.tlogger.WriteAddObjectInfo(name, info)
		}

		c.storage.AddObjectInfo(name, info)
	}

	if !c.hasPermission(userID, opts.Level) {
		return 0, constants.ErrForbidden
	}

	if !errNotFound {
		if info.Server != _mainStorage {
			return 0, constants.ErrServerNotFound
		}

		return _mainStorage, nil
	}

	if err := c.storage.CreateObject(name, opts); err != nil {
		return 0, fmt.Errorf("can't create object: %w", err)
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteCreateObject(name)
	}

	info = models.ObjectInfo{
		Server: _mainStorage,
		Level:  opts.Level,
	}

	c.storage.AddObjectInfo(name, info)

	return _mainStorage, nil
}

func (c *Core) GetFromObject(ctx context.Context, userID int, object, key string, opts models.GetFromObjectOptions) (v string, err error) {
	return v, c.withContext(ctx, func() error {
		v, err = c.getFromObject(ctx, userID, object, key, opts)
		return err
	})
}

func (c *Core) getFromObject(ctx context.Context, userID int, object, key string, opts models.GetFromObjectOptions) (string, error) {
	info, err := c.getObjectInfo(object)
	if err != nil {
		return "", err
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return "", constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		serverNumber := toServerNumber(opts.Server)

		cl, ok := c.servers.GetServerByID(serverNumber)
		if !ok || cl == nil {
			return "", constants.ErrServerNotFound
		}

		resp, err := cl.GetFromObject(ctx, object, key, opts)
		if err != nil {
			return "", err
		}

		return resp.Value, nil
	}

	if !c.hasPermission(userID, pkg.SafeDeref(opts.Level)) {
		return "", constants.ErrForbidden
	}

	v, err := c.storage.GetFromObject(object, key)
	if err != nil {
		return "", err
	}

	return v, nil
}

func (c *Core) SetToObject(ctx context.Context, userID int, object, key, val string, opts models.SetToObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.setToObject(ctx, userID, object, key, val, opts)
		return err
	})
}

func (c *Core) setToObject(ctx context.Context, userID int, object, key, val string, opts models.SetToObjectOptions) (int32, error) {
	info, err := c.storage.GetObjectInfo(object)
	if err != nil {
		return 0, fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return 0, constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return 0, constants.ErrServerNotFound
		}

		opts.Server = nil

		err = cl.SetToObject(ctx, object, key, val, opts)
		if err != nil {
			return 0, err
		}

		return info.Server, nil
	}

	if info.Server != _mainStorage {
		return 0, constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return 0, constants.ErrForbidden
	}

	err = c.storage.SetToObject(object, key, val, opts)
	if err != nil {
		return 0, err
	}

	if c.cfg.TransactionLogger.On {
		c.tlogger.WriteSetToObject(object, key, val)
	}

	return _mainStorage, nil
}

func (c *Core) ObjectToJSON(ctx context.Context, userID int, name string, opts models.ObjectToJSONOptions) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	info, err := c.storage.GetObjectInfo(name)
	if err != nil {
		return "", fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return "", constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return "", constants.ErrServerNotFound
		}

		res, err := cl.ObjectToJSON(ctx, name, opts)
		if err != nil {
			return "", err
		}

		return res.Object, nil
	}

	if info.Server != _mainStorage {
		return "", constants.ErrServerNotFound
	}

	if !c.hasPermission(userID, info.Level) {
		return "", constants.ErrForbidden
	}

	objJSON, err := c.storage.ObjectToJSON(name)
	if err != nil {
		return "", err
	}

	return objJSON, nil
}

func (c *Core) IsObject(ctx context.Context, userID int, name string, opts models.IsObjectOptions) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	_, err := c.storage.GetObjectInfo(name)
	if err != nil {
		return false, fmt.Errorf("can't get object info: %w", err)
	}

	return true, nil
}

func (c *Core) Size(ctx context.Context, userID int, name string, opts models.SizeOptions) (size uint64, err error) {
	return size, c.withContext(ctx, func() error {
		size, err = c.size(ctx, userID, name, opts)
		return err
	})
}

func (c *Core) size(ctx context.Context, userID int, name string, opts models.SizeOptions) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	info, err := c.storage.GetObjectInfo(name)
	if err != nil {
		return 0, fmt.Errorf("can't get object info: %w", err)
	}

	if c.earlyObjectNotFound(opts.Server, info.Server) {
		return 0, constants.ErrServerNotFound
	}

	if !c.useMainStorage(opts.Server) {
		cl, ok := c.servers.GetServerByID(info.Server)
		if !ok || cl == nil {
			return 0, constants.ErrServerNotFound
		}

		res, err := cl.Size(ctx, name, opts)
		if err != nil {
			return 0, err
		}

		return res.Size, nil
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

func (c *Core) DeleteObject(ctx context.Context, userID int, name string, opts models.DeleteObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, userID, name, opts)
	})
}

func (c *Core) deleteObject(ctx context.Context, userID int, name string, opts models.DeleteObjectOptions) error {
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

		err = cl.DeleteObject(ctx, name, opts)
		if err != nil {
			return err
		}

		c.storage.DeleteObjectInfo(name)

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

func (c *Core) AttachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.attachToObject(ctx, userID, dst, src, opts)
	})
}

func (c *Core) attachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) error {
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

		err = cl.AttachToObject(ctx, dst, src, opts)
		if err != nil {
			return err
		}

		return nil
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

func (c *Core) DeleteAttr(ctx context.Context, userID int, key string, object string, opts models.DeleteAttrOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, userID, key, object, opts)
	})
}

func (c *Core) deleteAttr(ctx context.Context, userID int, key, object string, opts models.DeleteAttrOptions) error {
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

		err = cl.DeleteAttr(ctx, key, object, opts)
		if err != nil {
			return err
		}

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
