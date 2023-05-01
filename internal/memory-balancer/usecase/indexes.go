package usecase

import (
	"context"
	"fmt"
)

var ErrNoActiveClients = fmt.Errorf("error while creating area: no clients")

func (uc *UseCase) Index(ctx context.Context, name string) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	if a, ok := uc.indexes[name]; ok {
		return a, nil
	}

	cl, ok := uc.servers.GetClient()
	if !ok || cl == nil {
		return 0, ErrNoActiveClients
	}

	err := cl.NewIndex(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("error while creating index: %w", err)
	}

	uc.indexes[name] = cl.Number
	return cl.Number, nil
}

func (uc *UseCase) GetFromIndex(ctx context.Context, index, key string, serverNumber int32) (string, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	var ok bool
	if serverNumber, ok = uc.indexes[index]; !ok {
		return "", fmt.Errorf("index %s does not exist", index)
	}

	cl, ok := uc.servers.GetClientByID(serverNumber)
	if !ok || cl == nil {
		return "", fmt.Errorf("no such server for index %s", index)
	}

	resp, err := cl.GetFromIndex(ctx, index, key)
	if err != nil {
		return "", err
	}

	cl.Total = resp.Total
	cl.Available = resp.Available
	return resp.Value, nil
}

func (uc *UseCase) SetToIndex(ctx context.Context, index, key, val string, uniques bool) (int32, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	var ok bool
	var serverNumber int32
	if serverNumber, ok = uc.indexes[index]; !ok {
		return serverNumber, fmt.Errorf("index %s does not exist", index)
	}

	cl, ok := uc.servers.GetClientByID(serverNumber)
	if !ok || cl == nil {
		return 0, fmt.Errorf("no such server for index %s", index)
	}

	resp, err := cl.SetToIndex(ctx, index, key, val, uniques)
	if err != nil {
		return 0, err
	}

	cl.Total = resp.Total
	cl.Available = resp.Available
	return serverNumber, nil
}

func (uc *UseCase) GetIndex(ctx context.Context, name string) (map[string]string, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var num int32
	var ok bool
	if num, ok = uc.indexes[name]; !ok {
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

	cl.Total = res.Total
	cl.Available = res.Available
	return res.Index, nil
}

func (uc *UseCase) IsIndex(ctx context.Context, name string) (bool, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	_, ok := uc.indexes[name]
	return ok, nil
}

func (uc *UseCase) Size(ctx context.Context, name string) (uint64, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	var num int32
	var ok bool
	if num, ok = uc.indexes[name]; !ok {
		return 0, fmt.Errorf("index %s does not exist", name)
	}

	cl, ok := uc.servers.GetClientByID(num)
	if !ok || cl == nil {
		return 0, fmt.Errorf("no such server for index %s", name)
	}

	res, err := cl.Size(ctx, name)
	if err != nil {
		return 0, err
	}

	cl.Total = res.Total
	cl.Available = res.Available

	return res.Size, nil
}
