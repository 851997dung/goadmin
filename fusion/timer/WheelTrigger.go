package timer

type VTWheelTrigger interface {
	VTWheelTimer
	impl_OnActive()
}

type WheelTrigger struct {
	WheelTimer
	pointTime uint64
	vt        VTWheelTrigger
}

func (trigger *WheelTrigger) InitVT(vt VTWheelTrigger) {
	trigger.vt = vt
	trigger.WheelTimer.InitVT(vt)
}

func (trigger *WheelTrigger) Init(
	triggerInterval uint64, triggerPoint uint64, loopCount uint32) {
	trigger.WheelTimer.Init(triggerInterval, loopCount)
	trigger.pointTime = triggerPoint
}

func (trigger *WheelTrigger) GoNextActiveTime() {
	actualTickTime := trigger.GetActualTickTime()
	for actualTickTime >= trigger.point_time() {
		trigger.go_next_point_time()
	}
}

func (trigger *WheelTrigger) OnActive() {
	trigger.vt.impl_OnActive()
}

func (trigger *WheelTrigger) impl_OnPrepare() bool {
	trigger.GoNextActiveTime()
	trigger.SetNextActiveTime(trigger.point_time())
	return trigger.WheelTimer.impl_OnPrepare()
}

func (trigger *WheelTrigger) impl_OnActivate() {
	trigger.GoNextActiveTime()
	trigger.SetNextActiveTime(trigger.point_time())
	trigger.OnActive()
}

func (trigger *WheelTrigger) set_point_time(t uint64) {
	trigger.pointTime = t
}
func (trigger *WheelTrigger) go_previous_point_time() {
	trigger.pointTime -= trigger.active_interval()
}
func (trigger *WheelTrigger) go_next_point_time() {
	trigger.pointTime += trigger.active_interval()
}
func (trigger *WheelTrigger) point_time() uint64 {
	return trigger.pointTime
}
