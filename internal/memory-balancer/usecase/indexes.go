package usecase

import (
	"context"
	"fmt"
	"itisadb/internal/memory-balancer/servers"
)

var ErrNoActiveClients = fmt.Errorf("error while creating index: no clients")
var ErrIndexNotFound = fmt.Errorf("index not found")

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
		cl, ok = uc.servers.GetClientByID(num)
	} else {
		cl, ok = uc.servers.GetClient()
	}

	if !ok || cl == nil {
		return 0, ErrNoActiveClients
	}

	err := cl.NewIndex(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("error while creating index: %w", err)
	}

	num = cl.GetNumber()
	uc.indexes[name] = num
	return num, nil
}

// TODO: handle serverNumber
func (uc *UseCase) GetFromIndex(ctx context.Context, index, key string, serverNumber int32) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	uc.mu.RLock()
	num, ok := uc.indexes[index]
	uc.mu.RUnlock()

	if !ok {
		return "", ErrIndexNotFound
	}

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return "", fmt.Errorf("no such server for index %s", index)
	}

	resp, err := cl.GetFromIndex(ctx, index, key)
	if err != nil {
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

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return 0, fmt.Errorf("no such server for index %s", index)
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
		return nil, fmt.Errorf("index %s does not exist", name)
	}

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return nil, fmt.Errorf("no such server for index %s", name)
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

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return 0, fmt.Errorf("no such server for index %s", name)
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

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return fmt.Errorf("no such server for index %s", name)
	}

	err := cl.DeleteIndex(ctx, name)
	if err != nil {
		return fmt.Errorf("error while deleting index: %w", err)
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

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return fmt.Errorf("no such server for index %s", dst)
	}

	err := cl.AttachToIndex(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("error while attaching index: %w", err)
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

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return ErrIndexNotFound
	}

	err := cl.DeleteAttr(ctx, attr, index)
	if err != nil {
		return fmt.Errorf("error while deleting attr: %w", err)
	}
	return nil
}
