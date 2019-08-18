package service

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"unsafe"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
)

const (
	ChunkSize = 1024 * 1024 * 10
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
	Paused
	Done
	QueueFree
)

type DownlodMsg struct {
	Control byte
	ID      byte
	Offset  int64
}

// Down is Downloader service
type Down struct {
	fileName    string
	manager     *Manager
	currOffet   int64 // file offset
	nOffset     int64 // offset used in next read (could be from resend or file)
	state       byte
	cControl    chan byte
	cResend     chan int64
	out         chan []byte
	pendingAcks int
	resends     []int64
	connid      wire.UID
}

func newDown(fname string, m *Manager, out chan []byte, cid wire.UID) *Down {
	resend := make([]int64, 10)
	for i := range resend {
		resend[i] = -1
	}

	return &Down{
		fileName:    fname,
		manager:     m,
		currOffet:   0,
		nOffset:     -1,
		state:       ReadyToGo,
		cControl:    make(chan byte),
		cResend:     make(chan int64),
		out:         out,
		pendingAcks: 0,
		resends:     resend,
		connid:      cid,
	}
}

func (d *Down) Consume(b []byte) {
	dr := wire.DownloaderReq{}
	_, err := wire.Decode(b, dr)
	if err != nil {
		return
	}
	if dr.Control == Resend {
		d.cResend <- dr.Offset
		return
	}
	d.cControl <- dr.Control
}

func (d *Down) Run() {
	logg.Debug("Running dl")
	defer d.queueFree()

	f, err := os.Open(d.fileName)
	if err != nil {
		return
	}
	defer f.Close()

	//h := sha256.New()

	//logg.Debug("Header size")
	//logg.Debug(HeadSize)
	//logg.Debug(unsafe.Sizeof(DownlodMsg{}))
	d.state = Running

	reader := bufio.NewReader(f)

	buffers := [3][]byte{
		make([]byte, (ChunkSize + HeadSize)),
		make([]byte, (ChunkSize + HeadSize)),
		make([]byte, (ChunkSize + HeadSize)),
	}
	bufindex := 0

	for {
		wholepacket := buffers[bufindex]
		fileContent := wholepacket[:ChunkSize]

		exit := d.waitForsignal()
		if exit {
			return
		}

		d.nextOffset()
		d.nOffset, err = f.Seek(d.nOffset, 0)
		if err != nil {
			logg.Debug(err)
			logg.Debug("could not seek")
			return
		}

		n, err := reader.Read(fileContent)
		if err != nil {
			if err == io.EOF {
				d.packAndSend(fileContent[:n], true)
				if d.state == QueueFree {
					//should be rare
					return
				}
				d.state = Done
				//continue
			}
			logg.Debug(err)
			return
		}
		// emit here
		d.packAndSend(fileContent[:n], false)
		if bufindex == 2 {
			bufindex = 0
		} else {
			bufindex++
		}
	}
}

func (d *Down) waitForsignal() bool {

	// check if we have ack slot open
	// if don't wait just send another packet
	if d.state == QueueFree {
		return true
	}

	if d.pendingAcks < AckSlots {
		if d.state != Done {
			return false
		}
	}

	for {
		select {
		case c, ok := <-d.cControl:
			if !ok {
				return true
			}
			switch c {
			case Start:
				if d.state == ReadyToGo {
					d.state = Running
					return false
				}
			case Ack:
				if d.pendingAcks != 0 {
					d.pendingAcks--
					return false
				}
			case Pause:
				d.state = Paused
			case Resume:
				d.state = Running
				return false
			case Stop:
				d.state = Done
				return true
			}
		case r, ok := <-d.cResend:
			if !ok {
				return true
			}
			if d.pendingAcks != 0 {
				d.pendingAcks--
			}
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
			// execution reached here means
			// resend slots is full
			d.resends = append(d.resends, r)
			return false
		}

	}

}

func (d *Down) nextOffset() {
	for i, rs := range d.resends {
		if rs != -1 {
			d.nOffset = rs
			d.resends[i] = -1
			return
		}

	}
	d.nOffset = d.currOffet
}

func (d *Down) packAndSend(b []byte, last bool) {

	if d.currOffet == d.nOffset {

		d.currOffet = d.currOffet + int64(len(b))

	}

	b2 := make([]byte, 10)
	binary.BigEndian.PutUint64(b2, uint64(d.nOffset))
	if last {
		b2[8] = Finished
	} else {
		b2[8] = File
	}

	b2[9] = 0
	b = append(b, b2...)

	head := wire.Header{
		Connid:  d.connid,
		CmdType: wire.CMD_DOWNLOAD_SERVICE,
		Flow:    wire.AgentToUser,
	}
	b = wire.AttachHeader(&head, b)
	//d.pendingAcks++
	d.out <- b
}

func (d *Down) queueFree() {
	logg.Debug("Freeing downloader service")
	d.state = QueueFree
	d.manager.closeService(d)
}

func (d *Down) Close() {

	// is this excessive/redundend ¯\_(ツ)_/¯

	select {
	case _, ok := <-d.cControl:
		if ok {
			close(d.cControl)
		}
	default:
		close(d.cControl)
	}

	select {
	case _, ok := <-d.cResend:
		if ok {
			close(d.cResend)
		}
	default:
		close(d.cResend)
	}
}
