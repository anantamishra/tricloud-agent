package service

const ChunkSize = 512 * 8

type DownlodMsg struct {
	ControlSig byte
	Offset     uint64
	Data       []byte
}

/* Downloader service */
type Down struct {
	//pass
}

func (d *Down) Consume([]byte) {

}

func (d *Down) Run() {

}

func (d *Down) Close() {

}
