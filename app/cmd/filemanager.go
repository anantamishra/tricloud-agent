package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"fmt"

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

	if path == "." {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home
		}
	}
	logg.Debug(path)

	fss, _ := afero.ReadDir(afs, path)

	dirlistReply := &wire.ListDirReply{}
	dirlistReply.Path = path

	dirlistReply.ParentPath = filepath.Dir(path)

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

func FmAction(rawdata []byte, out chan []byte, ctx context.Context) {
	p := logg.Debug
	p("__ACTION__")
	fmaction := &wire.FmActionReq{}
	head, err := wire.Decode(rawdata, fmaction)
	if err != nil {
		p(err)
		return
	}

	var response wire.FmActionRes

	switch fmaction.Action {
	case "copy":
		response = actionCopy(fmaction)
	case "move":
		response = doMove(fmaction)
	case "rename":
		response = doRename(fmaction)
	case "mkdir":
		response = doMkdir(fmaction)
	case "delete":
		response = doDelete(fmaction)
	case "info":
		//pass
	case "compress":
		//pass
	case "hash":
		//pass
	}
	bytes, err := wire.Encode(head.Connid, wire.CMD_FM_ACTION, wire.AgentToUser, response)
	if err != nil {
		p("@fm action")
		p(err)
		return
	}
	out <- bytes
}

func actionCopy(req *wire.FmActionReq) wire.FmActionRes {
	resp := wire.FmActionRes{}
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

func doRename(req *wire.FmActionReq) wire.FmActionRes {
	// rename works only with one file
	resp := wire.FmActionRes{}
	resp["action"] = "rename"
	src := path.Join(req.Basepath, req.Targets[0])
	dest := path.Join(req.Basepath, req.Destination)

	var afs = afero.NewOsFs()

	err := afs.Rename(src, dest)
	if err != nil {
		resp["error"] = []string{err.Error()}
	} else {
		resp["error"] = []string{}
	}

	return resp
}

func doMove(req *wire.FmActionReq) wire.FmActionRes {
	// rename works only with one file
	resp := wire.FmActionRes{}
	resp["action"] = "move"

	errorStr := []string{}

	for _, target := range req.Targets {
		src := path.Join(req.Basepath, target)
		dest := path.Join(req.Destination, target)

		err := os.Rename(src, dest)

		if err != nil {
			errorStr = append(errorStr, err.Error())
		}
	}

	resp["error"] = errorStr

	return resp
}

func doMkdir(req *wire.FmActionReq) wire.FmActionRes {
	resp := wire.FmActionRes{}
	resp["action"] = "mkdir"
	if !(len(req.Targets) == 1) {
		resp["error"] = []string{"incorrectNoPath"}
		return resp
	}

	path := path.Join(req.Basepath, req.Targets[0])

	var afs = afero.NewOsFs()
	if b, _ := afero.Exists(afs, path); b {
		resp["error"] = []string{"alreadyExits"}
		return resp
	}

	err := afs.Mkdir(path, 0755)
	if err != nil {
		resp["error"] = []string{err.Error()}
		return resp
	}
	resp["error"] = []string{}

	return resp
}

func doDelete(req *wire.FmActionReq) wire.FmActionRes {

	fmt.Printf("%+v", req)

	resp:= wire.FmActionRes{}
	resp["action"] = "delete"

	var afs = afero.NewOsFs()

	for _, traget := range req.Targets {
		path := path.Join(req.Basepath, traget)

		isdir, err := afero.IsDir(afs, path)

		if err != nil {
			resp["error"] = []string{err.Error()}
			return resp
		}

		if isdir {
			err = afs.RemoveAll(path)
			if err != nil {
				resp["error"] = []string{err.Error()}
				return resp
			}
		} else {
			err = afs.Remove(path)
			if err != nil {
				resp["error"] = []string{err.Error()}
				return resp
			}
		}

	}

	resp["error"] = []string{}

	return resp
}
