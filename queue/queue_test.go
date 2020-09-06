package queue

import (
	"go-components/queue/slice"
	"go-components/queue/two_lock"
	"testing"
)

func TestQueue(t *testing.T) {
	queues := map[string]Queue{
		"lock-free queue":   NewLKQueue(),
		"two-lock queue":    two_lock.NewTwoLock(),
		"slice-base queued": slice.NewSliceQueue(0),
	}

	for name, q := range queues {
		t.Run(name, func(t *testing.T) {
			count := 100
			for i := 0; i < count; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < count; i++ {
				v := q.Dequeue()
				if v == nil {
					t.Fatalf("got a nil value")
				}
				if v.(int) != i {
					t.Fatalf("expect %d but got %v", i, v)
				}
			}
		})
	}
}
