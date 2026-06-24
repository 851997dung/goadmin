package timer

import (
	"container/list"

	"admin/fusion/base"
)

type VTWheelTimer interface {
	impl_OnPrepare() bool
	impl_OnActivate()
	impl_Dispose()
}

type WheelTimer struct {
	activeInterval  uint64
	activeTickCount uint64
	loopCount       uint32
	n1, n2          int
	itr             *list.Element
	mgr             *WheelTimerMgr
	vt              VTWheelTimer
}

func (timer *WheelTimer) InitVT(vt VTWheelTimer) {
	timer.vt = vt
}

func (timer *WheelTimer) Init(activeInterval uint64, loopCount uint32) {
	timer.activeInterval = activeInterval
	timer.loopCount = loopCount
}

func (timer *WheelTimer) RePush(nextActiveTime uint64) {
	if timer.mgr != nil {
		timer.mgr.RePush(timer, nextActiveTime)
	}
}

func (timer *WheelTimer) GetActualTickTime() uint64 {
	return timer.mgr.tickParticle * timer.mgr.actualTickCount
}

func (timer *WheelTimer) GetCurrentTickTime() uint64 {
	return timer.mgr.tickParticle * timer.mgr.tickCount
}

func (timer *WheelTimer) GetNextActiveTime() uint64 {
	return timer.mgr.tickParticle * timer.activeTickCount
}

func (timer *WheelTimer) SetNextActiveTime(nextActiveTime uint64) {
	if nextActiveTime == 0 {
		timer.activeTickCount = timer.mgr.actualTickCount + base.
			MaxUint64(timer.activeInterval/timer.mgr.tickParticle, 1) - 1
	} else {
		timer.activeTickCount = nextActiveTime / timer.mgr.tickParticle
	}
	if timer.activeTickCount < timer.mgr.tickCount {
		timer.activeTickCount = timer.mgr.tickCount
	}
}

func (timer *WheelTimer) DetachMgr() {
	if timer.mgr != nil && timer.itr != nil {
		timer.mgr.allTimers[timer.n1][timer.n2].Remove(timer.itr)
		timer.mgr = nil
	}
}

func (timer *WheelTimer) Destroy() {
	timer.Dispose()
}

func (timer *WheelTimer) OnPrepare() bool {
	return timer.vt.impl_OnPrepare()
}

func (timer *WheelTimer) OnActivate() {
	timer.vt.impl_OnActivate()
}

func (timer *WheelTimer) Dispose() {
	timer.vt.impl_Dispose()
}

func (timer *WheelTimer) impl_OnPrepare() bool {
	return true
}

func (timer *WheelTimer) impl_Dispose() {
	timer.DetachMgr()
}

func (timer *WheelTimer) active_interval() uint64 {
	return timer.activeInterval
}
