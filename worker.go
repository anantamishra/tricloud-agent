package main

import (
	"context"
	"encoding/json"
	"log"

	cmd "github.com/indrenicloud/tricloud-agent/commands"
	"github.com/indrenicloud/tricloud-server/core"
)

// Worker coroutine, it recives packet, decodes it and runs functions commandbuff
// bashed on command type
func Worker(ctx context.Context, In, Out chan []byte) {

	for {
		select {
		case _ = <-ctx.Done():
			return
		case inData := <-In:
			msg := core.MessageFormat{}
			_ = json.Unmarshal(inData, &msg)

			cmdfunc, ok := cmd.CommandBuffer[msg.CmdType]
			if !ok {
				log.Println("Command not implemented")
				break
			}
			go cmdfunc(&msg, Out)
		}
	}

}
