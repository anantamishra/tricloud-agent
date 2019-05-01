package worker

import (
	"context"
	"os"

	"github.com/indrenicloud/tricloud-agent/wire"

	"github.com/indrenicloud/tricloud-agent/app/cmd"
	"github.com/indrenicloud/tricloud-agent/app/logg"
)

// Worker coroutine, it recives packet, decodes it and runs functions commandbuff
// bashed on command type
func Worker(ctx context.Context, In, Out chan []byte) {

	for {
		select {
		case _ = <-ctx.Done():
			return
		case inData := <-In:
			header, _ := wire.GetHeader(inData)

			if header.CmdType == wire.CMD_EXIT {
				os.Exit(0)
			}

			logg.Log("processing server cmd")

			cmdFunc, ok := cmd.CommandBuffer[header.CmdType]
			if !ok {
				logg.Log("Command not implemented")
				break
			}
			go cmdFunc(inData, Out)
		}
	}

}
