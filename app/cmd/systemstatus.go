package cmd

import (
	"context"
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
func SystemStatus(rawdata []byte, out chan []byte, ctx context.Context) {
	if sysinforunning {
		logg.Log("already systemstatus running")
		return
	}

	s := &wire.SysStatCmd{}
	wire.Decode(rawdata, s)
	sysinforunning = true
	defer func() { sysinforunning = false }()

	counter := int64(0)
	for {

		logg.Log("counter:", counter)

		rb, err := wire.Encode(wire.UID(0),
			wire.CMD_SYSTEM_STAT,
			wire.BroadcastUsers,
			systemStatus(time.Duration(s.Interval)*time.Second))
		if err == nil {
			out <- rb
		}
		//time.Sleep(time.Duration(s.Interval) * time.Second)

		if s.Timeout != 0 {
			counter = counter + s.Interval
			if counter >= (s.Timeout * s.Interval) {
				logg.Log("Exiting status emitting func")
				return
			}
		}
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
