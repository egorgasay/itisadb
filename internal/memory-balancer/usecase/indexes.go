package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"itisadb/internal/memory-balancer/servers"
)

var ErrObjectNotFound = fmt.Errorf("object not found")
var ErrServerNotFound = fmt.Errorf("server not found")

func (uc *UseCase) Object(ctx context.Context, name string) (s int32, err error) {
	return s, uc.withContext(ctx, func() error {
		s, err = uc.object(ctx, name)
		return err
	})
}

func (uc *UseCase) object(ctx context.Context, name string) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	var ok bool
	var cl *servers.Server
	var num int32

	if num, ok = uc.objects[name]; ok {
		cl, ok = uc.servers.GetServerByID(num)
	} else {
		cl, ok = uc.servers.GetServer()
	}

	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	err := cl.NewObject(ctx, name)
	if err != nil {
		return 0, err
	}

	num = cl.GetNumber()
	err = uc.storage.SaveObjectLoc(ctx, name, num)
	if err != nil {
		uc.logger.Warn(fmt.Sprintf("error while saving object: %s", err.Error()))
	}

	uc.objects[name] = num
	return num, nil
}

func (uc *UseCase) GetFromObject(ctx context.Context, object, key string, serverNumber int32) (v string, err error) {
	return v, uc.withContext(ctx, func() error {
		v, err = uc.getFromObject(ctx, object, key, serverNumber)
		return err
	})
}

func (uc *UseCase) getFromObject(ctx context.Context, object, key string, serverNumber int32) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[object]
	uc.mu.RUnlock()

	if !ok && serverNumber == 0 {
		return "", ErrObjectNotFound
	} else if serverNumber != 0 {
		num = serverNumber
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", ErrServerNotFound
	}

	resp, err := cl.GetFromObject(ctx, object, key)
	if err != nil {
		if errors.Is(err, servers.ErrNotFound) {
			return "", ErrNoData
		}
		return "", err
	}

	return resp.Value, nil
}

func (uc *UseCase) SetToObject(ctx context.Context, object, key, val string, uniques bool) (s int32, err error) {
	return s, uc.withContext(ctx, func() error {
		s, err = uc.setToObject(ctx, object, key, val, uniques)
		return err
	})
}

func (uc *UseCase) setToObject(ctx context.Context, object, key, val string, uniques bool) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[object]
	uc.mu.RUnlock()

	if !ok {
		return 0, ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	err := cl.SetToObject(ctx, object, key, val, uniques)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (uc *UseCase) ObjectToJSON(ctx context.Context, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[name]
	uc.mu.RUnlock()

	if !ok {
		return "", ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", ErrServerNotFound
	}

	res, err := cl.ObjectToJSON(ctx, name)
	if err != nil {
		return "", err
	}

	return res.Object, nil
}

func (uc *UseCase) IsObject(ctx context.Context, name string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	uc.mu.RLock()
	_, ok := uc.objects[name]
	uc.mu.RUnlock()

	return ok, nil
}

func (uc *UseCase) Size(ctx context.Context, name string) (size uint64, err error) {
	return size, uc.withContext(ctx, func() error {
		size, err = uc.size(ctx, name)
		return err
	})
}

func (uc *UseCase) size(ctx context.Context, name string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[name]
	uc.mu.RUnlock()

	if !ok {
		return 0, ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	res, err := cl.Size(ctx, name)
	if err != nil {
		return 0, err
	}

	return res.Size, nil
}

func (uc *UseCase) DeleteObject(ctx context.Context, name string) error {
	return uc.withContext(ctx, func() error {
		return uc.deleteObject(ctx, name)
	})
}

func (uc *UseCase) deleteObject(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[name]
	uc.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.DeleteObject(ctx, name)
	if err != nil {
		return err
	}

	uc.mu.Lock()
	delete(uc.objects, name)
	uc.mu.Unlock()

	return nil
}

func (uc *UseCase) AttachToObject(ctx context.Context, dst string, src string) error {
	return uc.withContext(ctx, func() error {
		return uc.attachToObject(ctx, dst, src)
	})
}

func (uc *UseCase) attachToObject(ctx context.Context, dst string, src string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	uc.mu.RLock()
	num, ok := uc.objects[dst]
	uc.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.AttachToObject(ctx, dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) DeleteAttr(ctx context.Context, attr string, object string) error {
	return uc.withContext(ctx, func() error {
		return uc.deleteAttr(ctx, attr, object)
	})
}

func (uc *UseCase) deleteAttr(ctx context.Context, attr string, object string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.objects[object]
	uc.mu.RUnlock()

	if !ok {
		return ErrObjectNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.DeleteAttr(ctx, attr, object)
	if err != nil {
		return err
	}
	return nil
}
