package cmd

import (
	"encoding/json"
	"testing"

	"github.com/shirou/gopsutil/process"
)

func Testjson(t *testing.T) {
	p, _ := process.Processes()

	b, _ := json.Marshal(p)
	t.Log("Test")
	t.Log(string(b))
}
