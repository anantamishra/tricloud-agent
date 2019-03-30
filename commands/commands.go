package commands

import (
	"fmt"

	"github.com/indrenicloud/tricloud-server/broker/core"
	"github.com/shirou/gopsutil/mem"
)

func init() {
	registerCommands()
}

// CommandFunc  is a common signature for different command type
// they take different string input and gives string output
type CommandFunc func(args ...string) string

// CommandBuffer contain the mapping of different command type to their mapping (func)
var CommandBuffer map[core.CommandType]CommandFunc

// all commands will be registered from here
func registerCommands() {
	// internal commands
	CommandBuffer[core.CMD_SYSTEM_INFO] = SystemInfo
}

// SystemInfo gives system info of machine
func SystemInfo(args ...string) string {

	v, _ := mem.VirtualMemory()
	sysinfo := fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	return sysinfo
}
