package cmd

import (
	"context"

	"github.com/indrenicloud/tricloud-agent/wire"
)

func SysAction(rawdata []byte, out chan []byte, ctx context.Context) {

	sareq := &wire.SystemActionReq{}
	_, err := wire.Decode(rawdata, sareq)
	if err != nil {
		return
	}

	switch sareq.Action {
	case "shutdown":
		doPoweroff()
	case "reboot":
		doReboot()
	}

}
