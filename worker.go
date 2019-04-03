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
		case inData := <-In:

			go func() {

				msg := core.MessageFormat{}

				_ = json.Unmarshal(inData, &msg)

				cmdfunc, ok := cmd.CommandBuffer[msg.CmdType]
				if !ok {
					log.Println("Command not implemented")
					return
				}
				outMsg := core.MessageFormat{
					Receiver: msg.Receiver,
					CmdType:  msg.CmdType,
					Arguments: map[string]string{
						"result": cmdfunc(""),
					},
				}

				Out <- outMsg.GetBytes()

			}()
		case _ = <-ctx.Done():
			return
		}
	}

}
