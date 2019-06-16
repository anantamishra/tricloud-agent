package cmd

import (
	"os/exec"

	"github.com/indrenicloud/tricloud-agent/app/logg"
)

func doPoweroff() {
	//err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err := exec.Command("systemctl", "poweroff").Run(); err != nil {
		logg.Debug(err)
	}
}

func doReboot() {
	//err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err := exec.Command("systemctl", "reboot").Run(); err != nil {
		logg.Debug(err)
	}

}
