package service

import (
	"bufio"
	"io"
	"os"
	"unsafe"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
)

const (
	ChunkSize = 500
	AckSlots  = 3
	HeadSize  = int(unsafe.Sizeof(DownlodMsg{}) + unsafe.Sizeof(wire.Header{}))
)

const (
	Start byte = iota
	Ack
	Resend
	Pause
	Resume
	Stop
	// agent side control states
	File
	Folder
	Hash
	Finished

	// service states
	ReadyToGo
	Running
	WaitingAck
	Paused
	Done
	QueueFree
)

type ControlSignal struct {
	signal byte
	data   interface{}
}

type DownlodMsg struct {
	Control byte
	Offset  int64
	id      byte
}

// Down is Downloader service
type Down struct {
	fileName  string
	currOffet int64
	backRead  int64
	state     byte
	cControl  chan ControlSignal
	out       chan []byte
	acks      []int64
	resends   []int64
}

func newDown(fname string, out chan []byte) *Down {

	ackslice := make([]int64, AckSlots)
	for i := range ackslice {
		ackslice[i] = -1
	}

	return &Down{
		fileName:  fname,
		currOffet: 0,
		backRead:  -1,
		state:     ReadyToGo,
		cControl:  make(chan ControlSignal),
		out:       out,
		acks:      ackslice,
	}
}

func (d *Down) Consume([]byte) {

}

func (d *Down) Run() {
	defer d.queueFree()

	f, err := os.Open(d.fileName)
	if err != nil {
		return
	}
	defer f.Close()

	//h := sha256.New()

	logg.Debug("Header size")
	logg.Debug(HeadSize)

	reader := bufio.NewReader(f)

	for {
		wholepacket := make([]byte, (ChunkSize + HeadSize))
		fileContent := wholepacket[:ChunkSize]

		exit := d.waitForsignal()
		if exit {
			return
		}

		d.nextOffset()

		n, err := reader.Read(fileContent)
		if err != nil {
			if err == io.EOF {
				// emit here
				d.packAndSend(fileContent[:n])
				break
			}
			logg.Debug(err)
			os.Exit(1)
		}
		// emit here
		d.packAndSend(fileContent[:n])

	}
	// job finished send hash maybe
	//fmt.Printf("%x", h.Sum(nil))
}

func (d *Down) waitForsignal() bool {

	// check if we have ack slot open
	// if don't wait just send another packet
	for _, a := range d.acks {
		if a == -1 {
			return false
		}
	}

	for {
		select {
		case c := <-d.cControl:

			switch c.signal {
			case Start:
				return false
				//
			case Ack:
				ack, ok := c.data.(int64)
				if ok {
					for i, a := range d.acks {
						if a == ack {
							d.acks[i] = -1
							return false
						}
					}
				}
				//pass
			case Resend:
				_, ok := c.data.(int64)
				if ok {
					// add to resend if it doesnot exist
					// remove pending ack of that offset
					//
				}

				return false
			case Pause:
				//
			case Resume:
				return false
			case Stop:
				return true
			}
		}

	}

}

func (d *Down) nextOffset() {
	if len(d.resends) == 0 {
		// seek(d.offset)
		return
	}
	return
}

func (d *Down) packAndSend(b []byte) {
	b = pack(b)

}

func (d *Down) queueFree() {
	logg.Debug("Freeing downloader service")

}

func (d *Down) Close() {

}

func pack(b []byte) []byte {
	return b
}
