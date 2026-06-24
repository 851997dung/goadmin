package timer

import (
	"container/list"
	"log"
	"runtime/debug"

	"admin/fusion/base"
	. "admin/fusion/buildin"
)

var capacity = [...]int{256, 64, 64, 64, 64}

type WheelTimerMgr struct {
	tickParticle uint64

	tickCount   uint64
	pointerSlot []int
	allTimers   [][]list.List

	actualTickCount uint64
}

func NewMgr(particle, curtime uint64) *WheelTimerMgr {
	mgr := new(WheelTimerMgr)
	mgr.tickParticle = particle
	mgr.tickCount = curtime / particle
	mgr.actualTickCount = curtime / particle
	mgr.pointerSlot = make([]int, len(capacity))
	mgr.allTimers = make([][]list.List, len(capacity))
	for i, n := 0, len(capacity); i < n; i++ {
		mgr.allTimers[i] = make([]list.List, capacity[i])
	}
	return mgr
}

func (mgr *WheelTimerMgr) Accelerate(acceltime uint64) {
	accelTickCount := acceltime / mgr.tickParticle
	if accelTickCount > mgr.tickCount {
		accelTickCount = mgr.tickCount
	}
	for i, n := 0, len(capacity); i < n; i++ {
		for j, m := 0, capacity[i]; j < m; j++ {
			for itr := mgr.allTimers[i][j].Front(); itr != nil; itr = itr.Next() {
				itr.Value.(*WheelTimer).activeTickCount -= accelTickCount
			}
		}
	}
	mgr.tickCount -= accelTickCount
}

func (mgr *WheelTimerMgr) Update(curtime uint64) {
	mgr.actualTickCount = curtime / mgr.tickParticle
	for mgr.tickCount < mgr.actualTickCount {
		mgr.CascadeAndTick()
		mgr.Relocate(mgr.Activate())
	}
}

func (mgr *WheelTimerMgr) RePush(timer *WheelTimer, nextActiveTime uint64) {
	var pendingTimers list.List
	pendingTimers.PushBack(timer)
	mgr.allTimers[timer.n1][timer.n2].Remove(timer.itr)
	timer.SetNextActiveTime(nextActiveTime)
	mgr.Relocate(&pendingTimers)
}

func (mgr *WheelTimerMgr) Push(timer *WheelTimer, firstActiveTime uint64) {
	timer.mgr = mgr
	timer.SetNextActiveTime(firstActiveTime)
	if timer.OnPrepare() {
		var pendingTimers list.List
		pendingTimers.PushBack(timer)
		mgr.Relocate(&pendingTimers)
	} else {
		timer.Destroy()
	}
}

func (mgr *WheelTimerMgr) Pop(timer *WheelTimer) {
	timer.Destroy()
}

func (mgr *WheelTimerMgr) Clear() {
	for i, n := 0, len(capacity); i < n; i++ {
		for j, m := 0, capacity[i]; j < m; j++ {
			tlist := &mgr.allTimers[i][j]
			for tlist.Len() > 0 {
				tlist.Front().Value.(*WheelTimer).Destroy()
			}
		}
	}
}

func (mgr *WheelTimerMgr) CascadeAndTick() {
	if 1+mgr.pointerSlot[0] >= capacity[0] {
		var pendingTimers list.List
		for i, n := 1, len(capacity); i < n; i++ {
			linear := mgr.pointerSlot[i] + 1
			slot := IfInt(linear < capacity[i], linear, 0)
			pendingTimers.PushBackList(&mgr.allTimers[i][slot])
			mgr.allTimers[i][slot].Init()
			if slot != 0 {
				break
			}
		}
		mgr.Relocate(&pendingTimers)
	}

	for i, n := 0, len(capacity); i < n; i++ {
		mgr.pointerSlot[i]++
		if mgr.pointerSlot[i] >= capacity[i] {
			mgr.pointerSlot[i] = 0
		} else {
			break
		}
	}

	mgr.tickCount++
}

func (mgr *WheelTimerMgr) Activate() *list.List {
	activeTimes := &mgr.allTimers[0][mgr.pointerSlot[0]]
	if activeTimes.Len() > 0 {
		itr := activeTimes.Front()
		anchor := activeTimes.PushFront(nil)
		for {
			timer := itr.Value.(*WheelTimer)
			timer.SetNextActiveTime(0)
			base.SafeHandler(func() { timer.OnActivate() })()
			itr = anchor.Next()
			if itr != nil && itr.Value == timer {
				activeTimes.MoveAfter(anchor, itr)
				itr = anchor.Next()
				switch timer.loopCount {
				case 0:
				case 1:
					timer.Destroy()
				default:
					timer.loopCount--
				}
			}
			if itr == nil {
				break
			}
		}
		activeTimes.Remove(anchor)
	}
	return activeTimes
}

func (mgr *WheelTimerMgr) Relocate(pendingTimers *list.List) {
	for pendingTimers.Len() > 0 {
		itr := pendingTimers.Front()
		timer := itr.Value.(*WheelTimer)
		timer.mgr = nil
		if timer.activeTickCount >= mgr.tickCount {
			evaluateValue := int(timer.activeTickCount - mgr.tickCount)
			for i, n := 0, len(capacity); i < n; i++ {
				linear := mgr.pointerSlot[i] + evaluateValue + 1
				if evaluateValue >= capacity[i] {
					evaluateValue = linear/capacity[i] - 1
				} else {
					slot := linear % capacity[i]
					tlist := &mgr.allTimers[i][slot]
					tlist.PushBack(pendingTimers.Remove(itr))
					timer.n1 = i
					timer.n2 = slot
					timer.itr = tlist.Back()
					timer.mgr = mgr
					break
				}
			}
		}
		if timer.mgr == nil {
			log.Printf("Relocate wheel timer(%d,%d) error.\n",
				timer.activeTickCount, mgr.tickCount)
			debug.PrintStack()
			pendingTimers.Remove(itr)
			timer.Destroy()
		}
	}
}
