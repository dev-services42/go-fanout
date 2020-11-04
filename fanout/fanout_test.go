package fanout_test

import (
	"context"
	"fmt"
	"github.com/dev-services42/go-fanout/fanout"
	"github.com/stretchr/testify/require"
	"testing"
)

func benchBroadcast(b *testing.B, subscriptions int) {
	s := fanout.New()
	ctx := context.Background()

	for i := 0; i < subscriptions; i++ {
		ch := s.Subscribe(ctx, fanout.AllowAll)
		go func() {
			for range ch {
			}
		}()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Broadcast(i)
	}
}

func BenchmarkFanOut_Broadcast(b *testing.B) {
	presets := []int{1, 10, 100}

	for _, n := range presets {
		b.Run(fmt.Sprintf("%d subscribers", n), func(b *testing.B) {
			benchBroadcast(b, n)
		})
	}
}

func TestFanOut(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	s := fanout.New()

	ch1 := s.Subscribe(ctx, fanout.AllowAll)
	ch2 := s.Subscribe(ctx, fanout.AllowAll)

	payload := "example payload"
	s.Broadcast(payload)

	require.Equal(t, payload, (<-ch1).(string))
	require.Equal(t, payload, (<-ch2).(string))

	payload2 := "example payload 2"
	s.Broadcast(payload2)

	require.Equal(t, payload2, (<-ch1).(string))

	payload3 := "example payload 3"
	s.Broadcast(payload3)

	require.Equal(t, payload3, (<-ch1).(string))

	require.Equal(t, payload2, (<-ch2).(string))
	require.Equal(t, payload3, (<-ch2).(string))

	cancel()

	s.Wait()
}
