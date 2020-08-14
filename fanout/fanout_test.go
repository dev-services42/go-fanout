package fanout_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dev-services42/go-fanout/fanout"
)

func benchFanOut(b *testing.B, subscriptions int) {
	s := fanout.New()
	ctx := context.Background()

	for i := 0; i < subscriptions; i++ {
		ch, err := s.Subscribe(ctx, fanout.AllowAll)
		assert.NoError(b, err)
		go func() {
			for range ch {
			}
		}()
	}

	payload := "example payload"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Broadcast(payload)
	}
}

func BenchmarkFanOut_Broadcast(b *testing.B) {
	presets := []int{1, 10, 100, 1000, 10000}

	for _, n := range presets {
		b.Run(fmt.Sprintf("%d subscription", n), func(b *testing.B) {
			benchFanOut(b, n)
		})
	}
}

func TestFanOut(t *testing.T) {
	ctx := context.Background()

	s := fanout.New()

	ch1, err := s.Subscribe(ctx, fanout.AllowAll)
	require.NoError(t, err)

	ch2, err := s.Subscribe(ctx, fanout.AllowAll)
	require.NoError(t, err)

	payload := "example payload"
	s.Broadcast(payload)

	require.Equal(t, payload, (<-ch1).(string))
	require.Equal(t, payload, (<-ch2).(string))
}
