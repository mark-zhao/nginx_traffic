package common

import (
	"time"
)
type Reflector struct {
	resyncPeriod time.Duration
	clock Clock
}
var neverExitWatch <-chan time.Time = make(chan time.Time)
func (r *Reflector) ResyncChan() (<-chan time.Time, func() bool) {
	// 如果resyncPeriod说明就不用定时同步，返回的是永久超时的定时器
	if r.resyncPeriod == 0 {
		return neverExitWatch, func() bool { return false }
	}
	// The cleanup function is required: imagine the scenario where watches
	// always fail so we end up listing frequently. Then, if we don't
	// manually stop the timer, we could end up with many timers active
	// concurrently.
	// 构建定时器
	t := r.clock.NewTimer(r.resyncPeriod)
	return t.C(), t.Stop
}

func NewReflector(resyncPeriod time.Duration) *Reflector {
	r := &Reflector{
		resyncPeriod:  resyncPeriod,
		clock:         &RealClock{},
	}
	return r
}


