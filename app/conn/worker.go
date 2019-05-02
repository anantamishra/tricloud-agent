package conn

import (
	"context"
	"os"

	"github.com/indrenicloud/tricloud-agent/wire"

	"github.com/indrenicloud/tricloud-agent/app/cmd"
	"github.com/indrenicloud/tricloud-agent/app/logg"
)

// Worker coroutine, it recives packet, decodes it and runs functions commandbuff
// bashed on command type
func (c *Connection) Worker() {

	for {
		select {
		case _ = <-c.workerctx.Done():
			return
		case inData := <-c.In:
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
			newctx1, _ := context.WithCancel(c.workerctx)
			go cmdFunc(inData, c.Out, newctx1)
		}
	}

}
