package state

type State struct {
	next    *State
	value   interface{}
	hasNext chan struct{}
}

func New(value interface{}) *State {
	state := newState()
	state.value = value
	state.next = newState()
	return state
}

func newState() *State {
	return &State{
		next:    nil,
		value:   nil,
		hasNext: make(chan struct{}),
	}
}

func (s *State) WaitChange() <-chan struct{} {
	return s.hasNext
}

func (s *State) Next() *State {
	return s.next
}

func (s *State) Value() interface{} {
	return s.value
}

func (s *State) Set(value interface{}) {
	s.next.next = newState()
	s.next.value = value
	close(s.hasNext)
}

func (s *State) Clone() *State {
	state := newState()
	state.next = s.next
	state.value = s.value
	state.hasNext = s.hasNext
	return state
}
