package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"

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

type FmActionReq struct {
	Action      string
	Basepath    string
	Targets     []string
	Destination string
	Options     []string
}

/*
type FmActionRes struct {
	Action string
	Status bool
	Errors []string
	Outputs []string
} */
type FmActionRes map[string]interface{}

func FmAction(rawdata []byte, out chan []byte, ctx context.Context) {
	p := logg.Debug

	fmaction := &FmActionReq{}
	head, err := wire.Decode(rawdata, fmaction)
	if err != nil {
		p(err)
		return
	}

	var response FmActionRes

	switch fmaction.Action {
	case "copy":
		response = actionCopy(fmaction)
	case "move":
		//pass
	case "rename":
		//pass
	case "mkdir":
		//pass
	case "delete":
		//pass
	case "info":
		//pass
	case "compress":

	}
	bytes, err := wire.Encode(head.Connid, wire.CMD_FM_ACTION, wire.AgentToUser, response)
	if err != nil {
		p("@fm action")
		p(err)
		return
	}
	out <- bytes
}

func actionCopy(req *FmActionReq) FmActionRes {
	resp := FmActionRes{}
	resp["action"] = "copy"

	var afs = afero.NewOsFs()
	es := []error{}

	errorAppend := func(_err error) {
		es = append(es, _err)
	}

	for _, node := range req.Targets {
		src := path.Join(req.Basepath, node)
		dest := path.Join(req.Destination, node)

		if src = path.Clean("/" + src); src == "" {
			errorAppend(os.ErrNotExist)
			continue
		}

		if dest = path.Clean("/" + dest); dest == "" {
			errorAppend(os.ErrNotExist)
			continue
		}

		if dest == src {
			errorAppend(os.ErrInvalid)
			continue
		}

		dir, err2 := afero.IsDir(afs, src)
		if err2 != nil {
			logg.Debug(err2)
			errorAppend(err2)
			continue
		}

		var err3 error

		if dir {
			err3 = doCopyDir(afs, src, dest)
		} else {
			err3 = doCopyFile(afs, src, dest)

		}
		if err3 != nil {
			errorAppend(err3)
		}

	}

	myerrrs := []string{}
	if myerrrs == nil {
		panic("FUCK")
	}
	if es == nil {
		panic("FUCK2")
	}
	for _, e := range es {
		myerrrs = append(myerrrs, e.Error())
	}
	resp["error"] = myerrrs

	return resp
}

/*  CREDIT to https://github.com/filebrowser/filebrowser/tree/master/fileutils */

func doCopyFile(fs afero.Fs, source, dest string) error {
	src, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	// Makes the directory needed to create the dst
	// file.
	err = fs.MkdirAll(filepath.Dir(dest), 0666)
	if err != nil {
		return err
	}

	// Create the destination file.
	dst, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the contents of the file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Copy the mode if the user can't
	// open the file.
	info, err := fs.Stat(source)
	if err != nil {
		err = fs.Chmod(dest, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}

func doCopyDir(fs afero.Fs, source, dest string) error {
	// Get properties of source.
	srcinfo, err := fs.Stat(source)
	if err != nil {
		return err
	}

	// Create the destination directory.
	err = fs.MkdirAll(dest, srcinfo.Mode())
	if err != nil {
		return err
	}

	dir, _ := fs.Open(source)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	var errs []error

	for _, obj := range obs {
		fsource := source + "/" + obj.Name()
		fdest := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories, recursively.
			err = doCopyDir(fs, fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			// Perform the file copy.
			err = doCopyFile(fs, fsource, fdest)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}

	if errString != "" {
		return errors.New(errString)
	}
	return nil
}

/*
func Template(req *FmActionReq) FmActionRes {
	var resp FmActionRes
	resp["action"] = "templateType"

	return resp
}

*/
