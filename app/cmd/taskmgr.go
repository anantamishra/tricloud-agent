package cmd

import (
	"context"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/logg"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/shirou/gopsutil/process"
)

func Taskmanager(rawdata []byte, out chan []byte, ctx context.Context) {

	tcmd := &wire.TaskMgrCmd{}
	head, err := wire.Decode(rawdata, tcmd)

	if err != nil {
		logg.Log("invalid data")
	}

	var counter int
	for {

		if tcmd.Timeout != 0 {
			counter = counter + (tcmd.Interval * int(time.Second))
			if counter > tcmd.Timeout {
				return
			}
		}

		tdata := taskmanager()

		outbyte, err := wire.Encode(head.Connid, wire.CMD_TASKMGR, wire.BroadcastUsers, tdata)
		if err != nil {
			out <- outbyte
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
