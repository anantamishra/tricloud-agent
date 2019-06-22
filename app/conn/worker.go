package conn

import (
	"context"
	"os"

	"github.com/indrenicloud/tricloud-agent/app/cmd"
	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/app/service"
	"github.com/indrenicloud/tricloud-agent/wire"
)

// Worker coroutine, it recives packet, decodes it and runs functions commandbuff
// bashed on command type
func (c *Connection) Worker() {

	m := service.NewManager(c.Out)
	defer m.Close()

	for {
		select {
		case _ = <-c.workerctx.Done():
			return
		case inData := <-c.In:
			header, _ := wire.GetHeader(inData)

			if header.CmdType == wire.CMD_EXIT {
				os.Exit(0)
			}

			logg.Debug("processing server cmd")

			cmdFunc, ok := cmd.CommandBuffer[header.CmdType]
			if ok {
				newctx1, _ := context.WithCancel(c.workerctx)
				logg.Debug("Found")
				go cmdFunc(inData, c.Out, newctx1)
				continue
			}

			if header.CmdType == wire.CMD_START_SERVICE ||
				header.CmdType == wire.CMD_DOWNLOAD_SERVICE {

				go m.Consume(header, inData)

			}

		}
	}

}
