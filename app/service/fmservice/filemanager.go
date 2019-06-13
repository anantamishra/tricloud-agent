package fmservice

import "github.com/spf13/afero"

type FileManager struct {
	fs afero.Fs
}

func NewFileManager() {
	//var AppFs = afero.NewOsFs()
}

/*
startfm
list
cd

*/
