package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/service"
	"github.com/indrenicloud/tricloud-agent/wire"
)

func main() {
	out := make(chan []byte)

	go func() {
		tt := time.NewTimer(5 * time.Minute)

		for {
			select {
			case b := <-out:
				//print(string(b))
				_, b1 := wire.GetHeader(b)
				offset := len(b1) - 10
				protoByte := b1[offset:]
				dm := service.DownlodMsg{}
				dm.Offset = int64(binary.BigEndian.Uint64(protoByte[:8]))
				dm.Control = protoByte[8]
				dm.ID = protoByte[9]
				fmt.Printf("%+v", dm)
				//fmt.Println(string(b1[:offset]))

			case <-tt.C:
				print("exitting")
			}
		}
	}()

	m := service.NewManager(out)

	//researchdata.txt
	dm := wire.StartServiceReq{
		Options:     []string{"abc.txt"},
		ServiceType: byte(wire.CMD_DOWNLOAD_SERVICE),
	}

	b, err := wire.Encode(wire.UID(0),
		wire.CMD_START_SERVICE,
		wire.UserToAgent,
		dm)
	if err != nil {
		return
	}
	head, _ := wire.GetHeader(b)

	m.Consume(head, b)

	var wg sync.WaitGroup

	wg.Add(1)

	wg.Wait()
	m.Close()

}
