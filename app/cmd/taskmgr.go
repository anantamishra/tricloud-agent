package cmd

import (
	"context"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/logg"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/shirou/gopsutil/process"
)

var taskmanagerrunning bool

func Taskmanager(rawdata []byte, out chan []byte, ctx context.Context) {
	logg.Log("Taskmanger service received data")
	tcmd := &wire.TaskMgrCmd{}
	head, err := wire.Decode(rawdata, tcmd)
	if taskmanagerrunning {
		return
	}
	taskmanagerrunning = true
	if err != nil {
		logg.Log("invalid data")
		return
	}

	defer func() {}()

	counter := int64(0)
	for {

		tdata := taskmanager()

		outbyte, err := wire.Encode(head.Connid, wire.CMD_TASKMGR, wire.BroadcastUsers, tdata)
		if err == nil {
			out <- outbyte
		}

		if tcmd.Timeout != 0 {
			counter = counter + tcmd.Interval

			if counter >= (tcmd.Timeout * tcmd.Interval) {
				logg.Log("Exiting taskmanager service, timeout")
				return
			}
		}

		time.Sleep(time.Duration(tcmd.Interval) * time.Second)
	}

}

func taskmanager() *wire.TaskMgrData {
	pss, err := process.Processes()

	if err != nil {
		return nil
	}

	taskdata := &wire.TaskMgrData{}

	for _, p := range pss {
		pinfo := &wire.ProcessInfo{}
		pinfo.CPU, _ = p.CPUPercent()
		pinfo.Command, _ = p.Cmdline()
		pinfo.USER, _ = p.Username()
		pinfo.PID = p.Pid
		meminfo, err := p.MemoryInfo()
		if err == nil {
			pinfo.MEM = meminfo.RSS
		}
		t, err := p.Times()
		if err == nil {
			pinfo.UpTime = t.Total()
		}
		childs, err := p.Children()

		var childpids []int32
		if err == nil {
			for _, c := range childs {
				childpids = append(childpids, c.Pid)
			}
		}
		pinfo.ChildPIDS = childpids

		taskdata.Processes = append(taskdata.Processes, pinfo)
	}
	return taskdata
}
