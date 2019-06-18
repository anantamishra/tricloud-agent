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
	fileName    string
	currOffet   int64 // file offset
	nOffset     int64 // offset used in next read (could be from resend or file)
	state       byte
	cControl    chan ControlSignal
	out         chan []byte
	pendingAcks int
	resends     []int64
}

func newDown(fname string, out chan []byte) *Down {
	resend := make([]int64, 10)
	for i := range resend {
		resend[i] = -1
	}

	return &Down{
		fileName:    fname,
		currOffet:   -1,
		nOffset:     -1,
		state:       ReadyToGo,
		cControl:    make(chan ControlSignal),
		out:         out,
		pendingAcks: 0,
		resends:     resend,
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
			return
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
	if d.pendingAcks < AckSlots {
		return false
	}

	for {
		select {
		case c := <-d.cControl:

			switch c.signal {
			case Start:
				if d.state == ReadyToGo {
					d.state = Running
					return false
				}
			case Ack:
				if d.state == WaitingAck {
					d.pendingAcks--
					return false
				}

			case Resend:
				r, ok := c.data.(int64)
				if !ok {
					break
				}
				d.pendingAcks--
				for _, rs := range d.resends {
					if rs == r {
						//already pending
						return false
					}
				}
				for i, rs := range d.resends {
					if rs == -1 {
						d.resends[i] = r
						return false
					}
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
	for _, rs := range d.resends {
		if rs != -1 {
			d.nOffset = rs
			return
		}
	}
	d.nOffset = d.currOffet
}

func (d *Down) packAndSend(b []byte) {

	if d.currOffet == d.nOffset {

		d.currOffet = d.currOffet + int64(len(b))
	}

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
