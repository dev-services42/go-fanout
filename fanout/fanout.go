package fanout

import (
	"context"
	"github.com/dev-services42/go-fanout/fanout/state"
	"sync"
)

type FanOut struct {
	wg         *sync.WaitGroup
	stateMutex *sync.RWMutex
	state      *state.State
}

func New() *FanOut {
	return &FanOut{
		wg:         new(sync.WaitGroup),
		stateMutex: new(sync.RWMutex),
		state:      state.New(nil),
	}
}

func (s *FanOut) Subscribe(ctx context.Context, filterFn func(interface{}) bool) <-chan interface{} {
	ch := make(chan interface{})

	s.stateMutex.RLock()
	curState := s.state.Clone()
	s.stateMutex.RUnlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case <-curState.WaitChange():
				curState = curState.Next()
				value := curState.Value()
				if !filterFn(value) {
					continue
				}

				ch <- value
			}
		}
	}()

	return ch
}

func (s *FanOut) Broadcast(value interface{}) {
	s.stateMutex.Lock()
	s.state.Set(value)
	s.state = s.state.Next()
	s.stateMutex.Unlock()
}

func (s *FanOut) Wait() {
	s.wg.Wait()
}
