package channelqueue_test

import (
	"fmt"

	"github.com/delabroj/channelqueue"
)

func Example() {
	queue := channelqueue.MakeQueue(3)
	for i := 1; i <= 4; i++ {
		err := queue.Enqueue(i)
		if err != nil {
			fmt.Printf("failed to enqueue %v: queue full\n", i)
		}
	}
	for i := 1; i <= 4; i++ {
		v, err := queue.Dequeue()
		if err != nil {
			fmt.Println("failed to dequeue: queue empty")
			continue
		}
		fmt.Println(v)
	}
	// Output: failed to enqueue 4: queue full
	// 1
	// 2
	// 3
	// failed to dequeue: queue empty
}
