package app

import (
	"context"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/conn"
	"github.com/indrenicloud/tricloud-agent/app/worker"
)

var WAITTIME time.Duration = 2 * time.Second

func Run() {
	In := make(chan []byte)
	Out := make(chan []byte, 10)

	ErrorChannel := make(chan struct{})

	workerctx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go worker.Worker(workerctx, In, Out)

	for {

		//since both reader/writer can send error ...clearing just in case
		clearChannel(ErrorChannel)

		//new connection
		connctx, connCancel := context.WithCancel(context.Background())
		Conn := conn.NewConnection(connctx, In, Out, ErrorChannel)
		Conn.Run()

		select {
		case <-ErrorChannel:
			connCancel()

			//since both reader/writer can send error ...clearing just in case
			clearChannel(ErrorChannel)

			time.Sleep(WAITTIME)
		}

	}
}

func clearChannel(c chan struct{}) {
	select {
	case <-c:
	default:
	}
}
