package service

import (
	"bufio"
	"io"
	"os"
	"unsafe"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
)

const ChunkSize = 512

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
	fileName string
	state    byte
	cControl chan byte
}

func newDown() *Down {
	return &Down{}
}

func (d *Down) Consume([]byte) {

}

func (d *Down) Run() {
	defer d.queueFree(nil)

	f, err := os.Open(d.fileName)
	if err != nil {
		return
	}
	defer f.Close()

	//h := sha256.New()

	headsize := int(unsafe.Sizeof(DownlodMsg{}) + unsafe.Sizeof(wire.Header{}))
	logg.Debug("Header size")
	logg.Debug(headsize)

	reader := bufio.NewReader(f)

	for {
		wholepacket := make([]byte, (ChunkSize + headsize))
		fileContent := wholepacket[:ChunkSize]

		n, err := reader.Read(fileContent)
		if err != nil {
			if err == io.EOF {
				// emit here
				d.packAndSend(fileContent[:n+headsize])
				d.queueFree(nil)
				break
			}
			logg.Debug(err)
			os.Exit(1)
		}
		// emit here
		d.packAndSend(fileContent[:n+headsize])
		d.waitForsignal()
	}
	// job finished send hash maybe
	//fmt.Printf("%x", h.Sum(nil))
}

func (d *Down) waitForsignal() {
	for {
		select {
		case c := <-d.cControl:
			switch c {
			case Start:
				//
			case Ack:
				//pass
			case Resend:
				//pass
			case Pause:
				//pass
			case Resume:
				//pass
			case Stop:
				//pass
			}
		}

	}

}

func (d *Down) packAndSend(b []byte) {
	b = pack(b)

}

func (d *Down) queueFree(err error) {
	logg.Debug("Freeing downloader service")
	if err != nil {
		logg.Debug(err)
	}

}

func (d *Down) Close() {

}

func pack(b []byte) []byte {
	return b
}
