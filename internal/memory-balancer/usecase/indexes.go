package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"itisadb/internal/memory-balancer/servers"
)

var ErrIndexNotFound = fmt.Errorf("index not found")
var ErrServerNotFound = fmt.Errorf("server not found")

func (uc *UseCase) Index(ctx context.Context, name string) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	var ok bool
	var cl *servers.Server
	var num int32

	if num, ok = uc.indexes[name]; ok {
		cl, ok = uc.servers.GetServerByID(num)
	} else {
		cl, ok = uc.servers.GetServer()
	}

	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	err := cl.NewIndex(ctx, name)
	if err != nil {
		return 0, err
	}

	num = cl.GetNumber()
	err = uc.storage.SaveIndexLoc(ctx, name, num)
	if err != nil {
		uc.logger.Warn(fmt.Sprintf("error while saving index: %s", err.Error()))
	}

	uc.indexes[name] = num
	return num, nil
}

func (uc *UseCase) GetFromIndex(ctx context.Context, index, key string, serverNumber int32) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[index]
	uc.mu.RUnlock()

	if !ok && serverNumber == 0 {
		return "", ErrIndexNotFound
	} else if serverNumber != 0 {
		num = serverNumber
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return "", ErrServerNotFound
	}

	resp, err := cl.GetFromIndex(ctx, index, key)
	if err != nil {
		if errors.Is(err, servers.ErrNotFound) {
			return "", ErrNoData
		}
		return "", err
	}

	return resp.Value, nil
}

func (uc *UseCase) SetToIndex(ctx context.Context, index, key, val string, uniques bool) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[index]
	uc.mu.RUnlock()

	if !ok {
		return 0, ErrIndexNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return 0, ErrServerNotFound
	}

	err := cl.SetToIndex(ctx, index, key, val, uniques)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (uc *UseCase) GetIndex(ctx context.Context, name string) (map[string]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[name]
	uc.mu.RUnlock()

	if !ok {
		return nil, ErrIndexNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return nil, ErrServerNotFound
	}

	res, err := cl.GetIndex(ctx, name)
	if err != nil {
		return nil, err
	}

	return res.Index, nil
}

func (uc *UseCase) IsIndex(ctx context.Context, name string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	uc.mu.RLock()
	_, ok := uc.indexes[name]
	uc.mu.RUnlock()

	return ok, nil
}

func (uc *UseCase) Size(ctx context.Context, name string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[name]
	uc.mu.RUnlock()

	if !ok {
		return 0, ErrIndexNotFound
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

func (uc *UseCase) DeleteIndex(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[name]
	uc.mu.RUnlock()

	if !ok {
		return ErrIndexNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.DeleteIndex(ctx, name)
	if err != nil {
		return err
	}

	uc.mu.Lock()
	delete(uc.indexes, name)
	uc.mu.Unlock()

	return nil
}

func (uc *UseCase) AttachToIndex(ctx context.Context, dst string, src string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	uc.mu.RLock()
	num, ok := uc.indexes[dst]
	uc.mu.RUnlock()

	if !ok {
		return ErrIndexNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.AttachToIndex(ctx, dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) DeleteAttr(ctx context.Context, attr string, index string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[index]
	uc.mu.RUnlock()

	if !ok {
		return ErrIndexNotFound
	}

	cl, ok := uc.servers.GetServerByID(num)
	if !ok || cl == nil {
		return ErrServerNotFound
	}

	err := cl.DeleteAttr(ctx, attr, index)
	if err != nil {
		return err
	}
	return nil
}
