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
	ch, err := f.Subscribe(ctx, fanout.AllowAll)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Consumer %d is ready", id))
	for payload := range ch {
		fmt.Println(fmt.Sprintf("Consumer %d receive payload %v", id, payload))
	}
}

func main() {
	f := fanout.New()
	ctx := context.Background()

	go simpleConsumer(ctx, 1, f)
	go simpleConsumer(ctx, 2, f)
	go simpleConsumer(ctx, 3, f)

	for i := 0; i < 100; i++ {
		f.Broadcast(fmt.Sprintf("hello %d", i))
		time.Sleep(time.Second)
	}

	f.Wait()
}
```

