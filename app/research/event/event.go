package event

// IEvent is a event type interface
type IEvent interface {
	IsSingleton() bool
	String() string
}

var eventMap = map[uint16]string{
	0: "Default",
	1: "DiskLimit",
	2: "CPULimit",
	3: "MemoryLimit",
	4: "WatchProcess",
	5: "WatchFileSize",
	6: "WatchFile",
}

// Event type
type Event uint16

func (e Event) String() string {
	eventStr, ok := eventMap[uint16(e)]
	if ok {
		return eventStr
	}
	return ""
}

func (e Event) IsSingleton() bool {
	return false
}

type EventMessage struct {
	Name      string
	Type      uint16
	Message   string
	Timestamp int64
}

type EventCommand struct {
	Name   string
	Type   uint16
	Action string
	Args   []string
}

type EventCommandReply struct {
	Sucess string
}
