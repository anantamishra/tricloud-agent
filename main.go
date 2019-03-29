package main

import (
	"context"
	"time"
)

func main() {
	In := make(chan []byte)
	Out := make(chan []byte)
	//worker := make(chan []byte)
	ErrorChannel := make(chan struct{})

	workerctx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go Worker(workerctx, In, Out)

start:
	//since both reader/writer can send error ...clearing just in case
	select {
	case <-ErrorChannel:
	default:
	}

	connctx, connCancel := context.WithCancel(context.Background())
	Conn := NewConnection(connctx, In, Out, ErrorChannel)
	Conn.Run()

	select {
	case <-ErrorChannel:
		connCancel()

		//since both reader/writer can send error ...clearing just in case
		select {
		case <-ErrorChannel:
		default:
		}
		time.Sleep(2 * time.Second)
		goto start
	}

}
