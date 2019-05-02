package cmd

import (
	"context"

	"github.com/indrenicloud/tricloud-agent/wire"
)

// CommandFunc  is a common signature for different command type
// they take different string input and gives string output
type CommandFunc func(rawdata []byte, out chan []byte, ctx context.Context)

// CommandBuffer contain the mapping of different command type to their mapping (func)
var CommandBuffer map[wire.CommandType]CommandFunc

func init() {
	registerCommands()
}

// all commands will be registered from here
func registerCommands() {
	// internal commands

	CommandBuffer = map[wire.CommandType]CommandFunc{
		wire.CMD_SYSTEM_STAT: SystemStatus,
		wire.CMD_TERMINAL:    Terminal,
	}
}
