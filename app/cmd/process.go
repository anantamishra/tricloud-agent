package cmd

import (
	"context"
	"errors"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/shirou/gopsutil/process"
)

func ProcessAction(rawdata []byte, out chan []byte, ctx context.Context) {

	paCmd := &wire.ProcessActionCmd{}
	head, err := wire.Decode(rawdata, paCmd)
	if err != nil {
		logg.Log(err)
	}

	outbyte, err := wire.Encode(head.Connid,
		wire.CMD_PROCESS_ACTION,
		wire.AgentToUser,
		processAction(paCmd.PID, paCmd.Action),
	)
	if err != nil {
		out <- outbyte
	}
}

func processAction(pid int32, action string) *wire.ProcessActionData {
	paData := &wire.ProcessActionData{}
	switch action {
	case "kill":
		paData.Output = killProcess(pid).Error()
	case "pause":
		//pass
	}
	return paData
}

func killProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
	}
	return errors.New("Sucess")
}
