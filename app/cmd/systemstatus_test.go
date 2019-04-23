package cmd

import (
	"testing"
	"time"
)

func Test_systemStatus(t *testing.T) {
	t.Logf("%+v", systemStatus(time.Second*5))
}
