package fanout

import (
	"context"
	"github.com/dev-services42/go-fanout/fanout/state"
	"sync"
)

type FanOut struct {
	wg    *sync.WaitGroup
	state *state.State
}

func New() *FanOut {
	return &FanOut{
		wg:    new(sync.WaitGroup),
		state: state.New(),
	}
}

func (s *FanOut) Subscribe(ctx context.Context, filterFn func(interface{}) bool) (<-chan interface{}, error) {
	ch := make(chan interface{})

	s.wg.Add(1)
	ready := make(chan struct{})
	go func() {
		defer s.wg.Done()
		defer close(ch)

		close(ready)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				payload := s.state.Get()
				if !filterFn(payload) {
					continue
				}

				ch <- payload
			}
		}
	}()

	<-ready

	return ch, nil
}

func (s *FanOut) Broadcast(payload interface{}) {
	s.state.Set(payload)
}

func (s *FanOut) Wait() {
	s.wg.Wait()
}
