package owl

import "sync"

// State holds the previous state snap shot
type State struct {
	lock  sync.Mutex
	store map[string]SnapShot
}

func (s *State) set(v SnapShot) {
	s.lock.Lock()
	s.store[v.Path] = v
	s.lock.Unlock()
}

func (s *State) get(id string) *SnapShot {
	v, ok := s.store[id]
	if !ok {
		return nil
	}
	return &v
}
