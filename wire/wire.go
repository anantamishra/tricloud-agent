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
	AgentToUser FlowType = iota
	UserToAgent
	BroadcastUsers
)

const (
	CMD_SERVER_HELLO CommandType = iota
	CMD_SYSTEM_STAT
	CMD_TERMINAL
	CMD_TASKMGR
	CMD_PROCESS_ACTION
	CMD_LIST_SERVICES
	CMD_SERVICE_ACTION
	CMD_BUILTIN_MAX
)

type Header struct {
	Connid  UID
	CmdType CommandType
	Flow    FlowType
}

func NewHeader(connid UID, cmdtype CommandType, flow FlowType) *Header {
	return &Header{
		Connid:  connid,
		CmdType: cmdtype,
		Flow:    flow,
	}
}

func AttachHeader(header *Header, body []byte) []byte {

	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, uint16(header.Connid))
	b[2] = byte(header.CmdType)
	b[3] = byte(header.Flow)

	return append(body, b...)
}

func UpdateHeader(header *Header, body []byte) []byte {
	offset := len(body) - int(unsafe.Sizeof(Header{}))
	return AttachHeader(header, body[:offset])
}

func GetHeader(raw []byte) (*Header, []byte) {
	h := &Header{}
	offset := len(raw) - int(unsafe.Sizeof(Header{}))
	headerbytes := raw[offset:]

	h.Connid = UID(binary.LittleEndian.Uint16(headerbytes[:2]))
	h.CmdType = CommandType(headerbytes[2])
	h.Flow = FlowType(headerbytes[3])
	return h, raw[:offset]
}

func Encode(connid UID, cmdtype CommandType, flow FlowType, v interface{}) ([]byte, error) {
	bodybyte, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	h := NewHeader(connid, cmdtype, flow)
	return AttachHeader(h, bodybyte), nil
}

func Decode(raw []byte, out interface{}) (*Header, error) {

	h, bodyraw := GetHeader(raw)

	err := json.Unmarshal(bodyraw, out)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// structs which will represent new command/message format

type SysStatCmd struct {
	Interval int32 //duration in sec
	Timeout  int32 // 0 means no timeouts
}

type SysStatData struct {
	CPUPercent   []float64
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
	Processes []*ProcessInfo
}

type ProcessInfo struct {
	PID       int32
	USER      string
	CPU       float64
	MEM       uint64
	UpTime    float64
	Command   string
	ChildPIDS []int32
}

type ProcessActionCmd struct {
	PID    int32
	Action string
	ID     int32
}

type ProcessActionData struct {
	Output string
	ID     int32
}

type ListServicesCmd struct {
}

type ListServicesMsg struct {
	Services []*ServiceInfo
}

type ServiceInfo struct {
	Name string
}
