package cmd

import (
	"context"
	"syscall"
)


func SysAction(rawdata []byte, out chan []byte, ctx context.Context) {
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	
	
}