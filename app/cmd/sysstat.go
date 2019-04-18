package cmd

import (
	"log"
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/shirou/gopsutil/mem"
)

var sysinforunning bool

// SystemStatus gives system status of machine
func SystemStatus(rawdata []byte, out chan []byte) {
	if sysinforunning {
		log.Println("already systemstatus running")
		return
	}

	s := wire.SysStatCmd{}
	wire.Decode(rawdata, s)
	sysinforunning = true
	defer func() { sysinforunning = false }()

	var counter int
	for {
		if s.Timeout == 0 {
			continue
		} else {
			counter = counter + (s.Interval * int(time.Second))
			if counter > s.Timeout {
				return
			}
		}

		v, _ := mem.VirtualMemory()
		sysstat := &wire.SysStatData{
			AvailableMem: v.Free,
			TotalMem:     v.Total,
		}

		rb, err := wire.Encode(wire.UID(0), wire.CMD_SYSTEMSTAT, wire.AgentToServer, sysstat)
		if err == nil {
			out <- rb
		}
		time.Sleep(time.Duration(s.Interval) * time.Second)

	}

}

func init() {
	sysinforunning = false
}
