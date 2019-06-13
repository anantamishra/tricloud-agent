package fmservice

import (
	"sync"
)

// singleton pattern kinda

var fmInstance *FmService
var lxinstance sync.Mutex

type FmService struct {
	fms     map[string]*FileManager
	myepoch int
}

func (fms *FmService) epoch() int {
	return fms.myepoch
}

func (fms *FmService) close() {

}

func (fms *FmService) getFilemanger(sessionid string) *FileManager {
	return nil
}

func newFmService(epoch int) *FmService {
	return &FmService{

		fms: make(map[string]*FileManager),
	}
}

func GetFilemanager(sessionid string, epoch int) *FileManager {
	lxinstance.Lock()
	defer lxinstance.Unlock()
	if fmInstance != nil {
		if fmInstance.epoch() != epoch {
			fmInstance.close()

			// create new fmservice
			fmInstance = newFmService(epoch)
		}
	}

	return fmInstance.getFilemanger(sessionid)
}
