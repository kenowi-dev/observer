package observer

import (
	"sync"
)

type Observer interface {
	AddCallback(f func(int)) int
	Unregister(int)
	PeriodicAsync() RunningObserver
	Periodic()
	OneShotAsync() RunningObserver
	OneShot()
}

type RunningObserver interface {
	AddCallback(f func(int)) int
	Unregister(int)
	Stop()
}

type callback[T func(int) | func(int, bool)] struct {
	f  T
	id int
}

type observer struct {
	callbacks []*callback[func(int)]
	mutex     sync.Mutex
	nextId    int
	quit      chan bool
}

func newObserver() *observer {
	return &observer{
		callbacks: make([]*callback[func(int)], 0),
		mutex:     sync.Mutex{},
		nextId:    0,
		quit:      make(chan bool, 5),
	}
}

func (o *observer) AddCallback(f func(int)) int {
	o.mutex.Lock()
	c := callback[func(int)]{
		f:  f,
		id: o.nextId,
	}
	o.callbacks = append(o.callbacks, &c)
	o.nextId++
	o.mutex.Unlock()
	return c.id
}

func (o *observer) Unregister(id int) {
	o.mutex.Lock()

	idx := -1
	for i, c := range o.callbacks {
		if c.id == id {
			idx = i
		}
	}
	if idx != -1 {
		// if the id is found, replace the position with the last element and reduce the length by 1
		o.callbacks[idx] = o.callbacks[len(o.callbacks)-1]
		o.callbacks = o.callbacks[:len(o.callbacks)-1]
	}
	o.mutex.Unlock()
}

func (o *observer) Stop() {
	o.quit <- true
}

func (o *observer) allCallbacks() {
	o.mutex.Lock()
	for _, c := range o.callbacks {
		go c.f(c.id)
	}
	o.mutex.Unlock()
	return
}
