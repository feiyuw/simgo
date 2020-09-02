package server

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/feiyuw/simgo/protocols"
)

var (
	serverStorage = &servers{M: map[uint64]*Server{}}
	nextServerID  uint64
)

type Server struct {
	sync.RWMutex

	Id             uint64                 `json:"id"`
	Name           string                 `json:"name"`
	Protocol       string                 `json:"protocol"`
	Port           int                    `json:"port"`
	Options        map[string]interface{} `json:"options"`
	Clients        []string               `json:"clients"` // clients identifier
	RpcServer      protocols.RpcServer
	Messages       []*Message
	MethodHandlers map[string]*MethodHandler
}

type servers struct {
	sync.RWMutex

	M map[uint64]*Server
}

func (s *servers) Add(value *Server) (uint64, error) {
	s.Lock()
	defer s.Unlock()
	value.Id = atomic.AddUint64(&nextServerID, 1)
	s.M[value.Id] = value
	return value.Id, nil
}

func (s *servers) Remove(key uint64) error {
	s.RLock()
	server, exists := s.M[key]
	s.RUnlock()

	if exists {
		if server.RpcServer != nil {
			if err := server.RpcServer.Close(); err != nil {
				return err
			}
		}
		delete(s.M, key)
	}
	return nil
}

func (s *servers) FindAll() ([]*Server, error) {
	s.RLock()
	defer s.RUnlock()
	items := make([]*Server, len(s.M))
	idx := 0
	for _, v := range s.M {
		items[idx] = v
		idx++
	}
	// sort servers with ID order
	sort.Slice(items, func(idx1, idx2 int) bool {
		return items[idx1].Id < items[idx2].Id
	})

	return items, nil
}

func (s *servers) FindOne(key uint64) (*Server, error) {
	s.RLock()
	defer s.RUnlock()
	if v, exists := s.M[key]; exists {
		return v, nil
	}
	return nil, errors.New("not found")
}
