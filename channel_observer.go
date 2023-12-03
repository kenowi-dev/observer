package observer

type channelObs struct {
	*observer
	ch <-chan bool
}

func NewChanelObs(ch <-chan bool) Observer {
	return &channelObs{
		observer: newObserver(),
		ch:       ch,
	}
}

func (o *channelObs) PeriodicAsync() RunningObserver {
	go o.Periodic()
	return o
}

func (o *channelObs) Periodic() {
	for {
		o.OneShot()
	}
}

func (o *channelObs) OneShotAsync() RunningObserver {
	go o.OneShot()
	return o
}

func (o *channelObs) OneShot() {
	select {
	case <-o.quit:
		return
	case <-o.ch:
		o.allCallbacks()
	}
}
