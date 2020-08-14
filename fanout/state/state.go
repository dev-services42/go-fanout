package state

import "sync"

type State struct {
	cond    *sync.Cond
	message interface{}
}

func New() *State {
	return &State{
		cond:    sync.NewCond(new(sync.Mutex)),
		message: nil,
	}
}

func (s *State) Get() interface{} {
	defer s.cond.L.Unlock()

	s.cond.L.Lock()
	s.cond.Wait()

	return s.message
}

func (s *State) Set(message interface{}) {
	s.cond.L.Lock()
	s.message = message
	s.cond.L.Unlock()

	s.cond.Broadcast()
}
