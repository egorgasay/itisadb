package core

import (
	"context"
	"fmt"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

func (c *Core) Object(ctx context.Context, userID uint, name string, opts models.ObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.object(ctx, userID, name, opts)
		return err
	})
}

func (c *Core) object(ctx context.Context, userID uint, name string, opts models.ObjectOptions) (int32, error) {
	num, ok := c.objects[name]

	if c.useMainStorage(opts.Server) {
		if ok {
			if num != mainStorage {
				return 0, constants.ErrServerNotFound
			}
			return mainStorage, nil
		}

		if err := c.storage.CreateObject(name, opts); err != nil {
			return 0, fmt.Errorf("can't create object: %w", err)
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteCreateObject(name)
		}

		c.mu.Lock()
		defer c.mu.Unlock()
		c.objects[name] = mainStorage

		return mainStorage, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	err := cl.NewObject(ctx, name, opts)
	if err != nil {
		return 0, err
	}

	num = cl.GetNumber()

	if c.cfg.TransactionLoggerConfig.On {
		err = c.tlogger.SaveObjectLoc(ctx, name, num)
		if err != nil {
			c.logger.Warn(fmt.Sprintf("error while saving object: %s", err.Error()))
		}
	}

	c.objects[name] = num
	return num, nil
}

func (c *Core) GetFromObject(ctx context.Context, userID uint, object, key string, opts models.GetFromObjectOptions) (v string, err error) {
	return v, c.withContext(ctx, func() error {
		v, err = c.getFromObject(ctx, userID, object, key, opts)
		return err
	})
}

func (c *Core) getFromObject(ctx context.Context, userID uint, object, key string, opts models.GetFromObjectOptions) (string, error) {
	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if c.useMainStorage(opts.Server) {
		if ok && num != mainStorage {
			return "", constants.ErrServerNotFound
		}

		v, err := c.storage.GetFromObject(object, key)
		if err != nil {
			return "", err
		}

		return v, nil
	}

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

func (c *Core) SetToObject(ctx context.Context, userID uint, object, key, val string, opts models.SetToObjectOptions) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.setToObject(ctx, userID, object, key, val, opts)
		return err
	})
}

func (c *Core) setToObject(ctx context.Context, userID uint, object, key, val string, opts models.SetToObjectOptions) (int32, error) {
	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return 0, constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return 0, constants.ErrServerNotFound
		}

		err := c.storage.SetToObject(object, key, val, opts)
		if err != nil {
			return 0, err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteSetToObject(object, key, val)
		}

		return mainStorage, nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	opts.Server = nil

	err := cl.SetToObject(ctx, object, key, val, opts)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (c *Core) ObjectToJSON(ctx context.Context, userID uint, name string, opts models.ObjectToJSONOptions) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return "", constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return "", constants.ErrServerNotFound
		}

		objJSON, err := c.storage.ObjectToJSON(name)
		if err != nil {
			return "", err
		}

		return objJSON, nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	res, err := cl.ObjectToJSON(ctx, name, opts)
	if err != nil {
		return "", err
	}

	return res.Object, nil
}

func (c *Core) IsObject(ctx context.Context, userID uint, name string, opts models.IsObjectOptions) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	c.mu.RLock()
	_, ok := c.objects[name]
	c.mu.RUnlock()

	return ok, nil
}

func (c *Core) Size(ctx context.Context, userID uint, name string, opts models.SizeOptions) (size uint64, err error) {
	return size, c.withContext(ctx, func() error {
		size, err = c.size(ctx, userID, name, opts)
		return err
	})
}

func (c *Core) size(ctx context.Context, userID uint, name string, opts models.SizeOptions) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return 0, constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return 0, constants.ErrServerNotFound
		}

		size, err := c.storage.Size(name)
		if err != nil {
			return 0, err
		}

		return size, nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	res, err := cl.Size(ctx, name, opts)
	if err != nil {
		return 0, err
	}

	return res.Size, nil
}

func (c *Core) DeleteObject(ctx context.Context, userID uint, name string, opts models.DeleteObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, userID, name, opts)
	})
}

func (c *Core) deleteObject(ctx context.Context, userID uint, name string, opts models.DeleteObjectOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return constants.ErrServerNotFound
		}

		err := c.storage.DeleteObject(name)
		if err != nil {
			return err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteDeleteObject(name)
		}

		c.mu.Lock()
		delete(c.objects, name)
		c.mu.Unlock()

		return nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.DeleteObject(ctx, name, opts)
	if err != nil {
		return err
	}

	c.mu.Lock()
	delete(c.objects, name)
	c.mu.Unlock()

	return nil
}

func (c *Core) AttachToObject(ctx context.Context, userID uint, dst, src string, opts models.AttachToObjectOptions) error {
	return c.withContext(ctx, func() error {
		return c.attachToObject(ctx, userID, dst, src, opts)
	})
}

func (c *Core) attachToObject(ctx context.Context, userID uint, dst, src string, opts models.AttachToObjectOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mu.RLock()
	num, ok := c.objects[dst]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return constants.ErrServerNotFound
		}

		err := c.storage.AttachToObject(dst, src)
		if err != nil {
			return err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteAttach(dst, src)
		}

		return nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.AttachToObject(ctx, dst, src, opts)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) DeleteAttr(ctx context.Context, userID uint, key string, object string, opts models.DeleteAttrOptions) error {
	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, userID, key, object, opts)
	})
}

func (c *Core) deleteAttr(ctx context.Context, userID uint, key, object string, opts models.DeleteAttrOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	if c.useMainStorage(opts.Server) {
		if num != mainStorage {
			return constants.ErrServerNotFound
		}

		err := c.storage.DeleteAttr(object, key)
		if err != nil {
			return err
		}

		if c.cfg.TransactionLoggerConfig.On {
			c.tlogger.WriteDeleteAttr(object, key)
		}

		return nil
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.DeleteAttr(ctx, key, object, opts)
	if err != nil {
		return err
	}
	return nil
}
