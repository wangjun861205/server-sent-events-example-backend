package pkg

import (
	"sync"
)

type AddrKeeper struct {
	addrs map[string]chan<- string
	lock  sync.RWMutex
}

func NewAddrKeeper() *AddrKeeper {
	return &AddrKeeper{
		addrs: make(map[string]chan<- string),
		lock:  sync.RWMutex{},
	}
}

func (a *AddrKeeper) register(key string, ch chan<- string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.addrs[key] = ch
}

func (a *AddrKeeper) unregister(key string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.addrs, key)
}

func (a *AddrKeeper) getAddr(key string) chan<- string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.addrs[key]
}
