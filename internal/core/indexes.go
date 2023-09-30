package core

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	servers2 "itisadb/internal/servers"
)

var ErrObjectNotFound = fmt.Errorf("object not found")
var ErrServerNotFound = fmt.Errorf("server not found")

func (c *Core) Object(ctx context.Context, name string) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.object(ctx, name)
		return err
	})
}

func (c *Core) object(ctx context.Context, name string) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var ok bool
	var cl *servers2.Server
	var num int32

	if num, ok = c.objects[name]; ok {
		cl, ok = c.servers.GetServerByID(num)
	} else {
		cl, ok = c.servers.GetServer()
	}

	if !ok || cl == nil {
		return 0, ErrServerNotFound
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

func (c *Core) GetFromObject(ctx context.Context, object, key string, serverNumber int32) (v string, err error) {
	return v, c.withContext(ctx, func() error {
		v, err = c.getFromObject(ctx, object, key, serverNumber)
		return err
	})
}

func (c *Core) getFromObject(ctx context.Context, object, key string, serverNumber int32) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok && serverNumber == 0 {
		return "", ErrObjectNotFound
	} else if serverNumber != 0 {
		num = serverNumber
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", ErrServerNotFound
	}

	resp, err := cl.GetFromObject(ctx, object, key)
	if err != nil {
		if errors.Is(err, servers2.ErrNotFound) {
			return "", ErrNoData
		}
		return "", err
	}

	return resp.Value, nil
}

func (c *Core) SetToObject(ctx context.Context, object, key, val string, uniques bool) (s int32, err error) {
	return s, c.withContext(ctx, func() error {
		s, err = c.setToObject(ctx, object, key, val, uniques)
		return err
	})
}

func (c *Core) setToObject(ctx context.Context, object, key, val string, uniques bool) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return 0, ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	err := cl.SetToObject(ctx, object, key, val, uniques)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (c *Core) ObjectToJSON(ctx context.Context, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return "", ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", ErrServerNotFound
	}

	res, err := cl.ObjectToJSON(ctx, name)
	if err != nil {
		return "", err
	}

	return res.Object, nil
}

func (c *Core) IsObject(ctx context.Context, name string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	c.mu.RLock()
	_, ok := c.objects[name]
	c.mu.RUnlock()

	return ok, nil
}

func (c *Core) Size(ctx context.Context, name string) (size uint64, err error) {
	return size, c.withContext(ctx, func() error {
		size, err = c.size(ctx, name)
		return err
	})
}

func (c *Core) size(ctx context.Context, name string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return 0, ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	res, err := cl.Size(ctx, name)
	if err != nil {
		return 0, err
	}

	return res.Size, nil
}

func (c *Core) DeleteObject(ctx context.Context, name string) error {
	return c.withContext(ctx, func() error {
		return c.deleteObject(ctx, name)
	})
}

func (c *Core) deleteObject(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[name]
	c.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
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

func (c *Core) AttachToObject(ctx context.Context, dst string, src string) error {
	return c.withContext(ctx, func() error {
		return c.attachToObject(ctx, dst, src)
	})
}

func (c *Core) attachToObject(ctx context.Context, dst string, src string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mu.RLock()
	num, ok := c.objects[dst]
	c.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.AttachToObject(ctx, dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) DeleteAttr(ctx context.Context, attr string, object string) error {
	return c.withContext(ctx, func() error {
		return c.deleteAttr(ctx, attr, object)
	})
}

func (c *Core) deleteAttr(ctx context.Context, attr string, object string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	c.mu.RLock()
	num, ok := c.objects[object]
	c.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := c.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.DeleteAttr(ctx, attr, object)
	if err != nil {
		return err
	}
	return nil
}
