package cmd

import (
	"encoding/json"

	"github.com/shirou/gopsutil/host"
)

func GetSystemInfo() []byte {
	info, err := host.Info()
	if err != nil {
		return nil
	}

	bt, err := json.Marshal(info)
	if err != nil {
		return nil
	}

	return bt
}
