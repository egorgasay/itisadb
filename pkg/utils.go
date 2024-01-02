package pkg

import (
	"context"
	"github.com/egorgasay/gost"
	"github.com/shirou/gopsutil/mem"
	"itisadb/internal/models"
	"sync"
)

func IsTheSameArray[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	tmp := make(map[T]struct{})
	for _, el := range a {
		tmp[el] = struct{}{}
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
}

func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}

func SafeDeref[T any](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}

func WithContext(ctx context.Context, fn func() error, pool chan struct{}, onStop func()) (err error) {
	ch := make(chan struct{})

	defer onStop()

	once := sync.Once{}
	done := func() { close(ch) }

	pool <- struct{}{}
	go func() {
		err = fn()
		once.Do(done)
		<-pool
	}()

	select {
	case <-ch:
		return err
	case <-ctx.Done():
		once.Do(done)
		return ctx.Err()
	}
}

func CalcRAM() (res gost.Result[models.RAM]) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return res.ErrNewUnknown(err.Error()) // TODO: ???
	}

	used := vmStat.Used / (1024 * 1024)
	total := vmStat.Total / (1024 * 1024)

	return res.Ok(models.RAM{
		Total:     total,
		Available: total - used,
	})
}
