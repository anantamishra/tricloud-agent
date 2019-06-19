package service

import (
	"github.com/indrenicloud/tricloud-agent/wire"
)

func serviceBuilder(head *wire.Header,
	sm *Manager,
	req *wire.StartServiceReq, out chan []byte) Servicer {

	switch head.CmdType {
	case wire.CMD_DOWNLOAD_SERVICE:
		return newDown(req.Options[0], sm, out, head.Connid)
	case wire.CMD_UPLOAD_SERVICE:
		//pass
	}
	return nil
}
