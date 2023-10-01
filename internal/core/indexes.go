package core

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"itisadb/internal/constants"
	servers "itisadb/internal/servers"
)

func (c *Core) Object(ctx context.Context, server *int32, name string) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.object(ctx, server, name)
		return err
	})
}

func (c *Core) object(ctx context.Context, server *int32, name string) (int32, error) {
	num, ok := c.objects[name]

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.useMainStorage(server) {
		if ok {
			return mainStorage, nil
		}

		if err := c.storage.CreateObject(name); err != nil {
			return 0, fmt.Errorf("can't create object: %w", err)
		}

		return mainStorage, nil
	}

	var cl *servers.Server

	if ok {
		cl, ok = c.servers.GetServerByID(num)
	} else {
		cl, ok = c.servers.GetServer()
	}

	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	err := cl.NewObject(ctx, name)
	if err != nil {
		return 0, err
	}

	num = cl.GetNumber()
	err = c.storage.SaveObjectLoc(ctx, name, num)
	if err != nil {
		c.logger.Warn(fmt.Sprintf("error while saving object: %s", err.Error()))
	}

	c.objects[name] = num
	return num, nil
}

func (c *Core) GetFromObject(ctx context.Context, server *int32, object, key string) (v string, err error) {
	return v, c.withContext(ctx, func() error {
		v, err = c.getFromObject(ctx, server, object, key)
		return err
	})
}

func (c *Core) getFromObject(ctx context.Context, server *int32, object, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	serverNumber := toServerNumber(server)

	if !ok && serverNumber == 0 {
		return "", constants.ErrObjectNotFound
	} else if serverNumber != 0 {
		num = serverNumber
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	resp, err := cl.GetFromObject(ctx, object, key)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return "", constants.ErrNoData
		}
		return "", err
	}

	return resp.Value, nil
}

func (c *Core) SetToObject(ctx context.Context, server *int32, object, key, val string, uniques bool) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.setToObject(ctx, server, object, key, val, uniques)
		return err
	})
}

func (c *Core) setToObject(ctx context.Context, server *int32, object, key, val string, uniques bool) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return 0, constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	err := cl.SetToObject(ctx, object, key, val, uniques)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (c *Core) ObjectToJSON(ctx context.Context, server *int32, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return "", constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", constants.ErrServerNotFound
	}

	res, err := cl.ObjectToJSON(ctx, name)
	if err != nil {
		return "", err
	}

	return res.Object, nil
}

func (c *Core) IsObject(ctx context.Context, server *int32, name string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	c.mu.RLock()
	_, ok := c.objects[name]
	c.mu.RUnlock()

	return ok, nil
}

func (c *Core) Size(ctx context.Context, server *int32, name string) (size uint64, err error) {
	return size, c.withContext(ctx, func() error {
		size, err = c.size(ctx, server, name)
		return err
	})
}

func (c *Core) size(ctx context.Context, server *int32, name string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return 0, constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, constants.ErrServerNotFound
	}

	res, err := cl.Size(ctx, name)
	if err != nil {
		return 0, err
	}

	return res.Size, nil
}

func (c *Core) DeleteObject(ctx context.Context, server *int32, name string) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, server, name)
	})
}

func (c *Core) deleteObject(ctx context.Context, server *int32, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.DeleteObject(ctx, name)
	if err != nil {
		return err
	}

	c.mu.Lock()
	delete(c.objects, name)
	c.mu.Unlock()

	return nil
}

func (c *Core) AttachToObject(ctx context.Context, server *int32, dst string, src string) error {
	return c.withContext(ctx, func() error {
		return c.attachToObject(ctx, server, dst, src)
	})
}

func (c *Core) attachToObject(ctx context.Context, server *int32, dst string, src string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mu.RLock()
	num, ok := c.objects[dst]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.AttachToObject(ctx, dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) DeleteAttr(ctx context.Context, server *int32, attr string, object string) error {
	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, server, attr, object)
	})
}

func (c *Core) deleteAttr(ctx context.Context, server *int32, attr string, object string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return constants.ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return constants.ErrServerNotFound
	}

	err := cl.DeleteAttr(ctx, attr, object)
	if err != nil {
		return err
	}
	return nil
}
