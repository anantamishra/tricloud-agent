package cmd

import (
	"fmt"
	"testing"

	"github.com/indrenicloud/tricloud-agent/wire"
)

func TestList(t *testing.T) {
	p := fmt.Println

	ll := listDirectory("/")
	p(ll)
}

func TestCopy(t *testing.T) {

	ff := &wire.FmActionReq{
		Action:      "copy",
		Basepath:    "/home/bing/golang/two/src/github.com/indrenicloud/tricloud-agent/app/cmd/testfolder/a",
		Targets:     []string{"targetfile"},
		Destination: "/home/bing/golang/two/src/github.com/indrenicloud/tricloud-agent/app/cmd/testfolder/b",
	}

	ffres := actionCopy(ff)

	t.Logf("%+v", ffres)

	ff2 := &wire.FmActionReq{
		Action:      "copy",
		Basepath:    "/home/bing/golang/two/src/github.com/indrenicloud/tricloud-agent/app/cmd/testfolder/a",
		Targets:     []string{"c"},
		Destination: "/home/bing/golang/two/src/github.com/indrenicloud/tricloud-agent/app/cmd/testfolder/b",
	}

	ffres2 := actionCopy(ff2)

	t.Logf("%+v", ffres2)

}
