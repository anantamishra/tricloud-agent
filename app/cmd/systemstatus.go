package cmd

import (
	"time"

	"github.com/shirou/gopsutil/cpu"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/shirou/gopsutil/mem"
)

var sysinforunning bool

func init() {
	sysinforunning = false
}

// SystemStatus gives system status of machine
func SystemStatus(rawdata []byte, out chan []byte) {
	if sysinforunning {
		logg.Log("already systemstatus running")
		return
	}

	s := wire.SysStatCmd{}
	wire.Decode(rawdata, s)
	sysinforunning = true
	defer func() { sysinforunning = false }()

	var counter int32
	for {
		// 0 means no timeouts
		if s.Timeout != 0 {
			counter = counter + (s.Interval * int32(time.Second))
			if counter > s.Timeout {
				return
			}
		}
		rb, err := wire.Encode(wire.UID(0),
			wire.CMD_SYSTEM_STAT,
			wire.BroadcastUsers,
			systemStatus(time.Duration(s.Interval)))
		if err == nil {
			out <- rb
		}
		time.Sleep(time.Duration(s.Interval) * time.Second)
	}

}

func systemStatus(interval time.Duration) *wire.SysStatData {

	sysstat := &wire.SysStatData{}

	v, err := mem.VirtualMemory()

	if err == nil {
		sysstat.AvailableMem = v.Free
		sysstat.TotalMem = v.Total
	}

	percent, err := cpu.Percent(interval, false)

	if err == nil {
		sysstat.CPUPercent = percent
	}
	return sysstat
}
