package service

import (
	"sync"

	"github.com/indrenicloud/tricloud-agent/app/logg"
	"github.com/indrenicloud/tricloud-agent/wire"
)

type Manager struct {
	out       chan []byte
	lServices sync.Mutex
	services  map[wire.UID]map[wire.CommandType]Servicer
}

func NewManager(out chan []byte) *Manager {
	return &Manager{
		out:       out,
		lServices: sync.Mutex{},
		services:  make(map[wire.UID]map[wire.CommandType]Servicer),
	}
}

func (m *Manager) addService(h *wire.Header, req *wire.StartServiceReq) {
	m.lServices.Lock()
	defer m.lServices.Unlock()

	s := serviceBuilder(h, m, req, m.out)
	if s == nil {
		return
	}

	perUIDservices, ok := m.services[h.Connid]
	if !ok {
		perUIDservices = make(map[wire.CommandType]Servicer)
		m.services[h.Connid] = perUIDservices
	}
	perUIDservices[h.CmdType] = s
}

func (m *Manager) getService(h *wire.Header) Servicer {
	perUIDservices, ok := m.services[h.Connid]
	if !ok {
		return nil
	}

	service, ok := perUIDservices[h.CmdType]
	if ok {
		return service
	}

	return nil
}

func (m *Manager) Consume(h *wire.Header, data []byte) {

	if h.CmdType == wire.CMD_START_SERVICE {
		ssrq := &wire.StartServiceReq{}
		_, err := wire.Decode(data, ssrq)
		if err != nil {
			logg.Debug("unmarshal Error")
			return
		}
		m.addService(h, ssrq)
		return
	}

	s := m.getService(h)
	if s == nil {
		logg.Debug("Did not find service")
		return
	}

}

func (m *Manager) Close() {
	m.lServices.Lock()
	defer m.lServices.Unlock()
	for _, sc := range m.services {
		for _, s := range sc {
			s.Close()
		}
	}
}

// called by service itself
func (m *Manager) closeService(sr Servicer) {
	m.lServices.Lock()
	defer m.lServices.Unlock()
	for _, sc := range m.services {
		for key, s := range sc {
			if s == sr {
				delete(sc, key)
				sr.Close()
			}
		}
	}
}
