package main

import (
	"sync"
	"time"
)

type (
	lock struct {
		ID       string `json:"id"`
		LockedAt int64  `json:"locked_at"`
		// ExpiresAt
	}
	register struct {
		sync.RWMutex
		locks map[string]*lock
	}
)

var reg *register

func init() {
	reg = &register{
		locks: map[string]*lock{},
	}
}

func (r *register) set(k string) (*lock, bool) {
	r.Lock()
	defer r.Unlock()

	if lock, exist := r.locks[k]; exist {
		// Lock already acquired
		return lock, false
	}

	r.locks[k] = &lock{
		ID:       k,
		LockedAt: time.Now().UTC().Unix(),
	}

	// Lock acquired
	return r.locks[k], true
}

func (r *register) unset(k string) bool {
	r.Lock()
	defer r.Unlock()

	_, exist := r.locks[k]
	delete(r.locks, k)

	return exist
}

func (r *register) all() []*lock {
	r.RLock()
	defer r.RUnlock()

	locks := []*lock{}
	for _, lock := range r.locks {
		locks = append(locks, lock)
	}

	return locks
}

func (r *register) get(k string) (*lock, bool) {
	r.RLock()
	defer r.RUnlock()

	lock, exist := r.locks[k]
	return lock, exist
}
