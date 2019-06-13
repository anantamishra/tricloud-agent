package cmd

import (
	"context"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/spf13/afero"
)

type ListDirReq struct {
	Path    string
	Options []string
}

type ListDirReply struct {
	Path    string
	FSNodes []FSNode
}

type FSNode struct {
	Name string
	Type string
	Size int64
}

func ListDirectory(rawdata []byte, out chan []byte, ctx context.Context) {

	ld := &ListDirReq{}
	head, err := wire.Decode(rawdata, ld)
	if err != nil {
		logg.Debug("Could not decode listdir cmd")
		return
	}
	ldr := listDirectory(ld.Path)

	bytes, err := wire.Encode(head.Connid, head.CmdType, head.Flow, ldr)
	if err != nil {
		logg.Debug("Could not encode listdir reply")
		return
	}
	out <- bytes
}

func listDirectory(path string) *ListDirReply {
	var afs = afero.NewOsFs()

	fss, _ := afero.ReadDir(afs, path)

	dirlistReply := &ListDirReply{}
	dirlistReply.Path = path

	for _, fs := range fss {

		fsn := FSNode{}
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
