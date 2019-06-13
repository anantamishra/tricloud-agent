package cmd

import (
	"context"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/spf13/afero"
)

func ListDirectory(rawdata []byte, out chan []byte, ctx context.Context) {
	logg.Debug("listing dir called")

	ld := &wire.ListDirReq{}
	head, err := wire.Decode(rawdata, ld)
	if err != nil {
		logg.Debug("Could not decode listdir cmd")
		return
	}
	ldr := listDirectory(ld.Path)

	bytes, err := wire.Encode(head.Connid, head.CmdType, wire.AgentToUser, ldr)
	if err != nil {
		logg.Debug("Could not encode listdir reply")
		return
	}
	out <- bytes
	logg.Debug("Outed listof dir")
}

func listDirectory(path string) *wire.ListDirReply {
	var afs = afero.NewOsFs()

	fss, _ := afero.ReadDir(afs, path)

	dirlistReply := &wire.ListDirReply{}
	dirlistReply.Path = path

	for _, fs := range fss {

		fsn := wire.FSNode{}
		fsn.Name = fs.Name()

		if fs.IsDir() {
			fsn.Type = "dir"
			fsn.Size = fs.Size()
		} else {
			fsn.Type = "file"
			fsn.Size = fs.Size()
		}

		dirlistReply.FSNodes = append(dirlistReply.FSNodes, fsn)
	}
	return dirlistReply
}

type CopyFileReq struct {
	BasePath string
	Filename string

	BaseDest string
}

func CopyFile(rawdata []byte, out chan []byte, ctx context.Context) {

}
