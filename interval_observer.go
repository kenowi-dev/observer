package observer

import (
	"time"
)

type IntervalObserver interface {
	AddEachTickCallback(f func(int, bool)) int
}

type RunningIntervalObserver interface {
	AddEachTickCallback(f func(int, bool)) int
}

type intervalObs struct {
	*observer
	everyTickCallbacks []*callback[func(int, bool)]
	predicate          func() bool
	interval           time.Duration
}

func NewIntervalObs(predicate func() bool, interval time.Duration) Observer {
	return &intervalObs{
		observer:  newObserver(),
		predicate: predicate,
		interval:  interval,
	}
}

func (o *intervalObs) AddEachTickCallback(f func(int, bool)) int {
	o.mutex.Lock()
	c := callback[func(int, bool)]{
		f:  f,
		id: o.nextId,
	}
	o.everyTickCallbacks = append(o.everyTickCallbacks, &c)
	o.nextId++
	o.mutex.Unlock()
	return c.id
}

func (o *intervalObs) OneShotAsync() RunningObserver {
	go o.OneShot()
	return o
}

func (o *intervalObs) OneShot() {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()
	select {
	case <-o.quit:
		return
	case <-ticker.C:
		notify := o.predicate()
		if notify {
			o.allCallbacks()
		}
		o.allEveryTickCallbacks(notify)
	}
}

func (o *intervalObs) PeriodicAsync() RunningObserver {
	go o.Periodic()
	return o
}

func (o *intervalObs) Periodic() {
	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()
	for {
		select {
		case <-o.quit:
			return
		case <-ticker.C:
			notify := o.predicate()
			if notify {
				o.allCallbacks()
			}
			o.allEveryTickCallbacks(notify)
		}
	}
}

func (o *intervalObs) allEveryTickCallbacks(notify bool) {
	o.mutex.Lock()
	for _, c := range o.everyTickCallbacks {
		go c.f(c.id, notify)
	}
	o.mutex.Unlock()
	return
}
