package cmd

import (
	"context"
	"os/exec"

	"github.com/indrenicloud/tricloud-agent/wire"
)

func RunScript(rawdata []byte, out chan []byte, ctx context.Context) {
	screq := wire.ScriptReq{}
	head, err := wire.Decode(rawdata, screq)
	if err != nil {
		print("decode err")
		return
	}

	outReq := wire.ScriptRes{
		Response: runScript(screq.Code),
	}

	respbyte, err := wire.Encode(head.Connid, head.CmdType, wire.AgentToUser, outReq)
	if err != nil {
		return
	}
	out <- respbyte
}

func runScript(code string) string {

	cmd := exec.Command("/bin/sh", "-c", code)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(output)
}
