package service_test

import (
	"testing"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/service"
)

func TestFoo(t *testing.T) {
	//t.Error()
	out := make(chan []byte)

	go func() {
		tt := time.NewTimer(30 * time.Second)

		for {
			select {
			case b := <-out:
				print(string(b))
			case <-tt.C:
				print("exitting")
			}
		}
	}()

	m := service.NewManager(out)
	m.Close()
}
