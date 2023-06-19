package storage

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
	"sync"
)

type Storage struct {
	mu *sync.RWMutex
}

func New() (*Storage, error) {
	return &Storage{
		mu: &sync.RWMutex{},
	}, nil
}

// RestoreIndexes restores index names.
func (s *Storage) RestoreIndexes(ctx context.Context) (map[string]int32, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	err := os.MkdirAll(".indexes", 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to create indexes dir: %w", err)
	}

	entry := make(map[string]int32)

	dir, err := os.ReadDir(".indexes")
	if err != nil {
		return nil, fmt.Errorf("failed to read indexes dir: %w", err)
	}

	for _, fe := range dir {
		if fe.IsDir() {
			continue
		}

		f, err := os.OpenFile(fmt.Sprintf(".indexes/%s", fe.Name()), os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open indexes: %w", err)
		}

		for {
			var index string
			_, err = fmt.Fscanln(f, &index)
			if err != nil {
				break
			}

			num, err := strconv.Atoi(fe.Name())
			if err != nil {
				log.Errorf("failed to convert server number: %v", err)
			} else {
				entry[index] = int32(num)
			}
		}

		f.Close()
	}

	return entry, nil
}

// SaveIndexLoc saves index location.
func (s *Storage) SaveIndexLoc(ctx context.Context, index string, server int32) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(fmt.Sprintf(".indexes/%d", server), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to open indexes: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(index + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to indexes: %w", err)
	}

	return nil
}
