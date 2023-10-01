package models

import "github.com/pbnjay/memory"

type RAM struct {
	Total     uint64
	Available uint64
}

// Update outputs the current, total and OS memory being used.
func (r *RAM) Update() RAM {
	// TODO: do not call it every time
	return RAM{
		Total:     memory.TotalMemory() / 1024 / 1024,
		Available: memory.FreeMemory() / 1024 / 1024,
	}
}
