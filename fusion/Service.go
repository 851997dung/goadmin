package fusion

import (
	"admin/fusion/base"
	"admin/fusion/timer"
	"log"
	"time"
)

type ServiceBase struct {
	Name  string
	Tasks chan func()
	//RPCMgr   *RPCManager
	TimerMgr *timer.WheelTimerMgr
}

func (s *ServiceBase) Init(name string) {
	s.Name = name
	s.Tasks = make(chan func(), 1024)
	//s.RPCMgr = NewRPCManager(32)
	s.TimerMgr = timer.NewMgr(1, uint64(time.Now().Unix()))
}

func (s *ServiceBase) Tick() {
	//s.RPCMgr.OnTick()
	s.TimerMgr.Update(uint64(time.Now().Unix()))
}

func (s *ServiceBase) Start() {
	go base.SafeHandler(func() {
		log.Printf("start service `%s` successfully.", s.Name)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for !IsStopService {
			base.SafeHandler(func() {
				select {
				//case task := <-s.Tasks:
				//	task()
				//case task := <-s.RPCMgr.tasks:
				//	task()
				case <-ticker.C:
					s.Tick()
				case <-SigStopService:
				}
			})()
		}
	})()
}

//func (s *ServiceBase) AsyncBlockInvoke(task func()) {
//	c := make(chan bool)
//	s.Tasks <- func() {
//		task()
//		close(c)
//	}
//	<-c
//}
