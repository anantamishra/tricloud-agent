package commands

import (
	"fmt"
	"time"

	"github.com/indrenicloud/tricloud-server/core"
	"github.com/shirou/gopsutil/mem"
)

// SystemInfo gives system info of machine
func SystemInfo(msg *core.MessageFormat, out chan []byte) {

	for i := 0; i < 10; i++ {
		v, _ := mem.VirtualMemory()
		sysinfo := fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
		out <- ConstructMessage(msg.ReceiverConnid, msg.CmdType, []string{sysinfo})
		time.Sleep(5 * time.Second)
	}

}
