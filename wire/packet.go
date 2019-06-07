package wire


type Byteable interface {
	GetBytes()[]byte
}


type Packet struct {
	header Byteable
	body interface{}
}

func NewPacket( head Byteable , body interface{}) *Packet {
	return &Packet{
		header:head,
		body:body,
	}
}

func(p *Packet) Encode()[]byte {
	return nil
}

func(p *Packet) Decode(b []byte) {
	
}
