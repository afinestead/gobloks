package manager

import (
	"gobloks/internal/types"
	"sync"
)

type StatusFlags struct {
	status types.Flags
	mtx    *sync.Mutex
}

func (f *StatusFlags) Set(status types.Flags) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.status |= status
}

func (f *StatusFlags) Get() types.Flags {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.status
}

func (f *StatusFlags) Has(status types.Flags) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	return f.status&status != 0
}

func (f *StatusFlags) Clear(status types.Flags) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.status &= ^status
}
