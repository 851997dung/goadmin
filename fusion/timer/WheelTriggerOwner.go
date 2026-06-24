package timer

import (
	"log"
	"math"
)

type WheelTriggerOwner struct {
	WheelRoutineOwner
	triggers map[uint32]map[*wheelTrigger]struct{}
}

func (owner *WheelTriggerOwner) CreateTrigger(triggerCycle TriggerCycle,
	triggerPoint TriggerPoint, cb func(), t uint32, repeats uint32) {
	switch triggerCycle {
	case ByMonthlyTrigger:
		owner.GetCacheWheelTimerMgr().Push(&new_wheelTriggerByMonthly(
			owner, triggerPoint, cb, t, repeats).WheelTimer, 0)
	default:
		triggerInterval := GetTriggerInterval(triggerCycle)
		triggerPoint := CalcPreviousTriggerPointTime(triggerCycle, triggerPoint)
		owner.CreateTrigger4tp(triggerInterval, triggerPoint, cb, t, repeats)
	}
}

func (owner *WheelTriggerOwner) CreateTriggerX(triggerCycle TriggerCycle,
	triggerPoint TriggerPoint, cb func(), repeats uint32) {
	var t uint32 = math.MaxUint32
	owner.CreateTrigger(triggerCycle, triggerPoint, cb, t, repeats)
}

func (owner *WheelTriggerOwner) CreateTrigger4tp(triggerInterval uint64,
	triggerPoint uint64, cb func(), t uint32, repeats uint32) {
	owner.GetCacheWheelTimerMgr().Push(&new_wheelTrigger(
		owner, triggerInterval, triggerPoint, cb, t, repeats).WheelTimer, 0)
}

func (owner *WheelTriggerOwner) CreateTriggerX4tp(triggerInterval uint64,
	triggerPoint uint64, cb func(), repeats uint32) {
	var t uint32 = math.MaxUint32
	owner.CreateTrigger4tp(triggerInterval, triggerPoint, cb, t, repeats)
}

func (owner *WheelTriggerOwner) RemoveTriggers(t uint32) {
	if owner.triggers != nil {
		if triggers := owner.triggers[t]; triggers != nil {
			for trigger := range triggers {
				trigger.Destroy()
			}
		}
	}
}

func (owner *WheelTriggerOwner) RemoveAllTriggers() {
	if owner.triggers != nil {
		for _, triggers := range owner.triggers {
			for trigger := range triggers {
				trigger.Destroy()
			}
		}
	}
}

func (owner *WheelTriggerOwner) HasTriggers(t uint32) bool {
	if owner.triggers != nil {
		if triggers := owner.triggers[t]; triggers != nil {
			return len(triggers) > 0
		}
	}
	return false
}

func (owner *WheelTriggerOwner) HasAnyTriggers() bool {
	if owner.triggers != nil {
		return len(owner.triggers) > 0
	}
	return false
}

func (owner *WheelTriggerOwner) FindTrigger(t uint32) *WheelTrigger {
	if owner.triggers != nil {
		if triggers := owner.triggers[t]; triggers != nil {
			for trigger := range triggers {
				return &trigger.WheelTrigger
			}
		}
	}
	return nil
}

type wheelTrigger struct {
	WheelTrigger
	cb    func()
	t     uint32
	owner *WheelTriggerOwner
}

func new_wheelTrigger(owner *WheelTriggerOwner, triggerInterval uint64,
	triggerPoint uint64, cb func(), t uint32, repeats uint32) *wheelTrigger {
	trigger := new(wheelTrigger)
	trigger.InitVT(trigger)
	trigger.Init(owner, triggerInterval, triggerPoint, cb, t, repeats)
	return trigger
}

func (trigger *wheelTrigger) Init(owner *WheelTriggerOwner, triggerInterval uint64,
	triggerPoint uint64, cb func(), t uint32, repeats uint32) {
	trigger.WheelTrigger.Init(triggerInterval, triggerPoint, repeats)
	trigger.cb = cb
	trigger.t = t
	trigger.owner = owner
}

func (trigger *wheelTrigger) DetachOwner() {
	if trigger.owner != nil {
		if trigger.owner.triggers != nil {
			if triggers := trigger.owner.triggers[trigger.t]; triggers != nil {
				delete(triggers, trigger)
				if len(triggers) == 0 {
					delete(trigger.owner.triggers, trigger.t)
				}
			}
		}
		trigger.owner = nil
	}
}

func (trigger *wheelTrigger) impl_OnPrepare() bool {
	var triggers map[*wheelTrigger]struct{}
	if trigger.owner.triggers == nil {
		trigger.owner.triggers = make(map[uint32]map[*wheelTrigger]struct{})
	} else {
		triggers = trigger.owner.triggers[trigger.t]
	}
	if triggers == nil {
		triggers = make(map[*wheelTrigger]struct{})
		trigger.owner.triggers[trigger.t] = triggers
	}
	triggers[trigger] = struct{}{}
	return trigger.WheelTrigger.impl_OnPrepare()
}

func (trigger *wheelTrigger) impl_OnActive() {
	trigger.cb()
}

func (trigger *wheelTrigger) impl_Dispose() {
	trigger.DetachOwner()
	trigger.WheelTrigger.impl_Dispose()
}

type wheelTriggerByMonthly struct {
	wheelTrigger
	triggerPoint TriggerPoint
}

func new_wheelTriggerByMonthly(owner *WheelTriggerOwner, triggerPoint TriggerPoint,
	cb func(), t uint32, repeats uint32) *wheelTriggerByMonthly {
	trigger := new(wheelTriggerByMonthly)
	trigger.InitVT(trigger)
	trigger.Init(owner, MAX_CYCLE_MONTHLY,
		CalcNextTriggerPointTimeByMonthly(triggerPoint), cb, t, repeats)
	trigger.triggerPoint = triggerPoint
	return trigger
}

func (trigger *wheelTriggerByMonthly) impl_OnPrepare() bool {
	if trigger.GetActualTickTime() >= trigger.pointTime {
		log.Println("CreateTriggerByMonthly is invalid.")
		return false
	}
	return trigger.wheelTrigger.impl_OnPrepare()
}

func (trigger *wheelTriggerByMonthly) impl_OnActive() {
	trigger.set_point_time(CalcNextTriggerPointTimeByMonthly(trigger.triggerPoint))
	trigger.SetNextActiveTime(trigger.point_time())
	trigger.wheelTrigger.impl_OnActive()
}
