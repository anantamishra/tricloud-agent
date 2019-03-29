package commands

import (
	"encoding/json"
	"log"
)

// CommandType could be byte since we won't have more than 255 commands
type CommandType int

const (
	CMD_SERVER_HELLO CommandType = iota
	CMD_SERVER_MAX
	CMD_SYSTEM_INFO = CMD_SERVER_MAX + 1
	CMD_BUILTIN_MAX
)

// CMD_ALL_MAX future use :wink
var CMD_ALL_MAX CommandType = CMD_BUILTIN_MAX

// MessageFormat is core message format
type MessageFormat struct {
	Receiver  string
	CmdType   CommandType
	Arguments map[string]string
}

func (m *MessageFormat) GetBytes() []byte {

	outByte, err := json.Marshal(m)
	if err != nil {
		log.Fatal("could not marsal msg")
		return nil
	}

	return outByte
}
