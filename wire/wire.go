package wire

import (
	"encoding/binary"
	"encoding/json"
	"unsafe"
)

type FlowType byte
type CommandType byte
type UID uint16

const (
	AgentToServer FlowType = iota
	UserToServer
	AgentToUser
	UserToAgent
)

const (
	CMD_SERVER_HELLO CommandType = iota
	CMD_SYSTEMSTAT
	CMD_TERMINAL
	CMD_TASKMGR
	CMD_LISTSERVICES
	CMD_ACTIONSERVICE
	CMD_BUILTIN_MAX
)

type Header struct {
	Connid  UID
	CmdType CommandType
	Flow    FlowType
}

type SysStatCmd struct {
	Interval int //duration in msec
	Timeout  int // 0 means no timeouts
}

type SysStatData struct {
	CPUPercent   []uint8
	TotalMem     uint64
	AvailableMem uint64
}

type TermCmd struct {
	Command string // "" defaults to bash
	Args    []string
	EnvVars []string
}

type TermData struct {
	Data       string
	ResizeInfo string
}

type TaskMgrCmd struct {
	Interval int //duration in msec
	Timeout  int // 0 means default timeout will be used
}

type TaskMgrData struct {
	Uptime    int64
	AvgLoad   int
	Battery   uint8
	Processes []ProcessInfo
}

type ProcessInfo struct {
	PID       int
	ParentPID int
	USER      string
	CPU       int
	MEM       int
	UpTime    int64
	Command   string
}

func AttachHeader(connid UID, cmdtype CommandType, flow FlowType, body []byte) []byte {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, uint16(connid))
	b[2] = byte(cmdtype)
	b[3] = byte(flow)

	return append(body, b...)
}

func GetHeader(packet []byte) *Header {
	h := &Header{}

	h.Connid = UID(binary.LittleEndian.Uint16(packet[:2]))
	h.CmdType = CommandType(packet[2])
	h.Flow = FlowType(packet[3])
	return h
}

func Encode(connid UID, cmdtype CommandType, flow FlowType, v interface{}) ([]byte, error) {
	bodybyte, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return AttachHeader(connid, cmdtype, flow, bodybyte), nil
}

func Decode(raw []byte, out interface{}) (*Header, error) {

	offset := len(raw) - int(unsafe.Sizeof(Header{}))
	h := GetHeader(raw[offset:])

	err := json.Unmarshal(raw[:offset], out)
	if err != nil {
		return nil, err
	}
	return h, nil
}
