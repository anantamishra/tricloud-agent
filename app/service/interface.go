package service

type Servicer interface {
	Run()
	Consume([]byte)
	Close()
}
