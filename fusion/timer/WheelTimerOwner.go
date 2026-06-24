package timer

import "math"

type WheelTimerOwner struct {
	WheelRoutineOwner
	timers map[uint32]map[*wheelTimer]struct{}
}

func (owner *WheelTimerOwner) CreateTimer(
	cb func(), t uint32, interval uint64, repeats uint32, firstime uint64) {
	owner.GetCacheWheelTimerMgr().Push(
		&new_wheelTimer(owner, cb, t, interval, repeats).WheelTimer, firstime)
}

func (owner *WheelTimerOwner) CreateTimerX(
	cb func(), interval uint64, repeats uint32, firstime uint64) {
	var t uint32 = math.MaxUint32
	owner.CreateTimer(cb, t, interval, repeats, firstime)
}

func (owner *WheelTimerOwner) RemoveTimers(t uint32) {
	if owner.timers != nil {
		if timers := owner.timers[t]; timers != nil {
			for timer := range timers {
				timer.Destroy()
			}
		}
	}
}

func (owner *WheelTimerOwner) RemoveAllTimers() {
	if owner.timers != nil {
		for _, timers := range owner.timers {
			for timer := range timers {
				timer.Destroy()
			}
		}
	}
}

func (owner *WheelTimerOwner) HasTimers(t uint32) bool {
	if owner.timers != nil {
		if timers := owner.timers[t]; timers != nil {
			return len(timers) > 0
		}
	}
	return false
}

func (owner *WheelTimerOwner) HasAnyTimers() bool {
	if owner.timers != nil {
		return len(owner.timers) > 0
	}
	return false
}

func (owner *WheelTimerOwner) FindTimer(t uint32) *WheelTimer {
	if owner.timers != nil {
		if timers := owner.timers[t]; timers != nil {
			for timer := range timers {
				return &timer.WheelTimer
			}
		}
	}
	return nil
}

type wheelTimer struct {
	WheelTimer
	cb    func()
	t     uint32
	owner *WheelTimerOwner
}

func new_wheelTimer(owner *WheelTimerOwner,
	cb func(), t uint32, interval uint64, repeats uint32) *wheelTimer {
	timer := new(wheelTimer)
	timer.InitVT(timer)
	timer.Init(interval, repeats)
	timer.cb = cb
	timer.t = t
	timer.owner = owner
	return timer
}

func (timer *wheelTimer) DetachOwner() {
	if timer.owner != nil {
		if timer.owner.timers != nil {
			if timers := timer.owner.timers[timer.t]; timers != nil {
				delete(timers, timer)
				if len(timers) == 0 {
					delete(timer.owner.timers, timer.t)
				}
			}
		}
		timer.owner = nil
	}
}

func (timer *wheelTimer) impl_OnPrepare() bool {
	var timers map[*wheelTimer]struct{}
	if timer.owner.timers == nil {
		timer.owner.timers = make(map[uint32]map[*wheelTimer]struct{})
	} else {
		timers = timer.owner.timers[timer.t]
	}
	if timers == nil {
		timers = make(map[*wheelTimer]struct{})
		timer.owner.timers[timer.t] = timers
	}
	timers[timer] = struct{}{}
	return timer.WheelTimer.impl_OnPrepare()
}

func (timer *wheelTimer) impl_OnActivate() {
	timer.cb()
}

func (timer *wheelTimer) impl_Dispose() {
	timer.DetachOwner()
	timer.WheelTimer.impl_Dispose()
}
