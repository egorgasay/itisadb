package storage

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
)

// RestoreObjects restores object names.
func (s *Storage) RestoreObjects(ctx context.Context) (map[string]int32, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	err := os.MkdirAll(".objects", 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to create objects dir: %w", err)
	}

	entry := make(map[string]int32)

	dir, err := os.ReadDir(".objects")
	if err != nil {
		return nil, fmt.Errorf("failed to read objects dir: %w", err)
	}

	for _, fe := range dir {
		if fe.IsDir() {
			continue
		}

		f, err := os.OpenFile(fmt.Sprintf(".objects/%s", fe.Name()), os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open objects: %w", err)
		}

		for {
			var object string
			_, err = fmt.Fscanln(f, &object)
			if err != nil {
				break
			}

			num, err := strconv.Atoi(fe.Name())
			if err != nil {
				log.Errorf("failed to convert server number: %v", err)
			} else {
				entry[object] = int32(num)
			}
		}

		f.Close()
	}

	return entry, nil
}

// SaveObjectLoc saves object location.
func (s *Storage) SaveObjectLoc(ctx context.Context, object string, server int32) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(fmt.Sprintf(".objects/%d", server), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to open objects: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(object + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to objects: %w", err)
	}

	return nil
}
