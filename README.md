# Fan-out

## Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/dev-services42/go-fanout/fanout"
	"time"
)

func simpleConsumer(ctx context.Context, id int, f *fanout.FanOut) {
	ch := f.Subscribe(ctx, fanout.AllowAll)

	fmt.Println(fmt.Sprintf("Consumer %d is ready", id))
	go func() {
		for payload := range ch {
			fmt.Println(fmt.Sprintf("Consumer %d receive payload `%v`", id, payload))
		}
	}()
}

func main() {
	f := fanout.New()
	ctx, cancel := context.WithCancel(context.Background())

	simpleConsumer(ctx, 1, f)
	simpleConsumer(ctx, 2, f)
	simpleConsumer(ctx, 3, f)

	for i := 0; i < 1; i++ {
		f.Broadcast(fmt.Sprintf("hello %d", i))
		time.Sleep(time.Second)
	}

	cancel()
	f.Wait()
}
```

