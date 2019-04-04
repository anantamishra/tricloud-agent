package commands

import (
	"github.com/indrenicloud/tricloud-server/core"
)

// CommandFunc  is a common signature for different command type
// they take different string input and gives string output
type CommandFunc func(msg *core.MessageFormat, out chan []byte)

// CommandBuffer contain the mapping of different command type to their mapping (func)
var CommandBuffer map[core.CommandType]CommandFunc

func init() {
	registerCommands()
}

// all commands will be registered from here
func registerCommands() {
	// internal commands

	CommandBuffer = map[core.CommandType]CommandFunc{
		core.CMD_SYSTEM_INFO: SystemInfo,
	}
}

func ConstructMessage(connid core.UID, cmdtype core.CommandType, result []string) []byte {

	resultpacket := core.MessageFormat{
		ReceiverConnid: connid,
		CmdType:        cmdtype,
		Results:        result,
	}
	return resultpacket.GetBytes()
}
