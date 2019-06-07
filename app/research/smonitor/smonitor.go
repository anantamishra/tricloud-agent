package smonitor

import (
	"context"
)

type Service func(context.Context, chan []byte, bool)

type ServiceInfo struct {
	TimeOut  int64
	Interval int64
	Entry    string // incase of terminal it could be bash
	Args     []string
}

// Smonitor is a system monitor service
type Smonitor struct {
	ctx             context.Context
	outChan         chan interface{}
	servicesMap     map[string]Service // terminal, sysmon, watchfile
	runningServices map[string]context.CancelFunc

	broadcasting bool
}

func NewSmonitor(ctx context.Context) *Smonitor {
	return &Smonitor{}
}

func (sm *Smonitor) StartService(ctx context.Context, name string, si *ServiceInfo) {

	_, ok := sm.runningServices[name]
	if ok {
		print("Service already started")
		return
	}

	switch name {
	case "terminal":
		//start terminal
		//ctx, ctxFunc = context.WithCancel(sm.ctx)
	}

}

func (sm *Smonitor) Out() <-chan interface{} {
	return sm.outChan
}

func (sm *Smonitor) passiveMonitor(ctx context.Context) {
	//pass
}
