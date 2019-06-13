package terminal

import (
	"context"
	"sync"
)

// singleton pattern kinda

var termInstance *TermService
var lxinstance sync.Mutex

type TermService struct {
	ctx context.Context
}

func NewTermService(ctx context.Context) *TermService {
	lxinstance.Lock()
	defer lxinstance.Unlock()
	if termInstance != nil {
		return termInstance
	}

	termInstance = new(TermService)
	termInstance.ctx = ctx

	return termInstance
}
