package service

import (
	"github.com/indrenicloud/tricloud-agent/wire"
)

func serviceBuilder(head *wire.Header,
	sm *Manager,
	req *wire.StartServiceReq) Servicer {

	switch head.CmdType {
	case wire.CMD_DOWNLOAD_SERVICE:
		//pass
	case wire.CMD_UPLOAD_SERVICE:
		//pass
	}
	return nil
}
