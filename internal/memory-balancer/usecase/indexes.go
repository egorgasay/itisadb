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

	return uc.Get(ctx, index+".__"+key, serverNumber)
}

func (uc *UseCase) SetToIndex(ctx context.Context, index, key, val string, serverNumber int32, uniques bool) (int32, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	var ok bool
	if serverNumber, ok = uc.indexes[index]; !ok {
		return serverNumber, fmt.Errorf("index %s does not exist", index)
	}

	return uc.Set(ctx, index+".__"+key, val, serverNumber, uniques)
}
