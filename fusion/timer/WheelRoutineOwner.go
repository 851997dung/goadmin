package timer

type VTWheelRoutineOwner interface {
	Impl_GetWheelTimerMgr() *WheelTimerMgr
}

type WheelRoutineOwner struct {
	mgr *WheelTimerMgr
	vt  VTWheelRoutineOwner
}

func (owner *WheelRoutineOwner) InitVT(vt VTWheelRoutineOwner) {
	owner.vt = vt
}

func (owner *WheelRoutineOwner) GetCacheWheelTimerMgr() *WheelTimerMgr {
	if owner.mgr != nil {
		return owner.mgr
	}
	owner.mgr = owner.GetWheelTimerMgr()
	return owner.mgr
}

func (owner *WheelRoutineOwner) GetWheelTimerMgr() *WheelTimerMgr {
	return owner.vt.Impl_GetWheelTimerMgr()
}

type WheelRoutineType struct {
	routineType uint32
}

func (t *WheelRoutineType) Init() {
	t.routineType = 65535
}

func (t *WheelRoutineType) NewUniqueRoutineType() uint32 {
	t.routineType++
	return t.routineType
}
