package timer

import (
	"time"
)

const MAX_CYCLE_MONTHLY = 60 * 60 * 24 * 31

type TriggerCycle int

const (
	ByMonthlyTrigger TriggerCycle = iota
	ByWeeklyTrigger
	ByDailyTrigger
	ByHourlyTrigger
	By5MinutelyTrigger
	ByMinutelyTrigger
)

type TriggerPoint struct {
	Wday int8
	Hour uint8
	Min  uint8
	Sec  uint8
}

func GetTriggerInterval(tc TriggerCycle) uint64 {
	switch tc {
	case ByWeeklyTrigger:
		return 60 * 60 * 24 * 7
	case ByDailyTrigger:
		return 60 * 60 * 24
	case ByHourlyTrigger:
		return 60 * 60
	case By5MinutelyTrigger:
		return 60 * 5
	case ByMinutelyTrigger:
		return 60
	}
	return 0
}

func GetTriggerPointTime(tc TriggerCycle, tp TriggerPoint) uint64 {
	t := time.Now()
	pt := t.Unix()
	switch tc {
	case ByWeeklyTrigger:
		pt += int64(int(tp.Wday)-int(t.Weekday())) * (60 * 60 * 24)
		fallthrough
	case ByDailyTrigger:
		pt += int64(int(tp.Hour)-t.Hour()) * (60 * 60)
		fallthrough
	case ByHourlyTrigger:
		pt += int64(int(tp.Min)-t.Minute()) * (60)
		fallthrough
	case ByMinutelyTrigger:
		pt += int64(int(tp.Sec)-t.Second()) * (1)
	}
	return uint64(pt)
}

func CalcPreviousTriggerPointTime(tc TriggerCycle, tp TriggerPoint) uint64 {
	pt := GetTriggerPointTime(tc, tp)
	if pt > uint64(time.Now().Unix()) {
		pt -= GetTriggerInterval(tc)
	}
	return pt
}

func CalcNextTriggerPointTime(tc TriggerCycle, tp TriggerPoint) uint64 {
	pt := GetTriggerPointTime(tc, tp)
	if pt <= uint64(time.Now().Unix()) {
		pt += GetTriggerInterval(tc)
	}
	return pt
}

func CalcPreviousTriggerPointTimeByMonthly(tp TriggerPoint) uint64 {
	t := time.Now()
	pt := time.Date(t.Year(), t.Month(),
		int(tp.Wday), int(tp.Hour), int(tp.Min), int(tp.Sec), 0, time.Local)
	if pt.After(t) {
		pt = pt.AddDate(0, -1, 0)
	}
	return uint64(pt.Unix())
}

func CalcNextTriggerPointTimeByMonthly(tp TriggerPoint) uint64 {
	t := time.Now()
	pt := time.Date(t.Year(), t.Month(),
		int(tp.Wday), int(tp.Hour), int(tp.Min), int(tp.Sec), 0, time.Local)
	if pt.Before(t) || pt.Equal(t) {
		pt = pt.AddDate(0, 1, 0)
	}
	return uint64(pt.Unix())
}
