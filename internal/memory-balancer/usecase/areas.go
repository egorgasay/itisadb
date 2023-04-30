package usecase

import (
	"context"
	"fmt"
)

var ErrNoActiveClients = fmt.Errorf("error while creating area: no clients")

func (uc *UseCase) NewArea(ctx context.Context, name string) (int32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	if a, ok := uc.areas[name]; ok {
		return a, nil
	}

	cl, ok := uc.servers.GetClient()
	if !ok || cl == nil {
		return 0, ErrNoActiveClients
	}

	uc.areas[name] = cl.Number
	return cl.Number, nil
}

func (uc *UseCase) GetFromArea(ctx context.Context, area, key string, serverNumber int32) (string, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	var ok bool
	if serverNumber, ok = uc.areas[area]; !ok {
		return "", fmt.Errorf("area %s does not exist", area)
	}

	return uc.Get(ctx, area+".__"+key, serverNumber)
}

func (uc *UseCase) SetToArea(ctx context.Context, area, key, val string, serverNumber int32, uniques bool) (int32, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	var ok bool
	if serverNumber, ok = uc.areas[area]; !ok {
		return serverNumber, fmt.Errorf("area %s does not exist", area)
	}

	return uc.Set(ctx, area+".__"+key, val, serverNumber, uniques)
}
