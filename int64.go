package channelqueue

import "sync"

// Int64Queue holds the current channel and the maximum length of the channel.
type Int64Queue struct {
	front   chan int64
	maxSize int
	mux     sync.Mutex
}

// MakeInt64Queue returns a pointer to an initialized queue struct.
func MakeInt64Queue(n int) *Int64Queue {
	ch := make(chan int64, 1)
	return &Int64Queue{front: ch, maxSize: n}
}

// Enqueue adds the value to the queue's channel. If the channel is full, the
// channel is replaced with another channel twice as large until maxSize is reached.
func (q *Int64Queue) Enqueue(n int64) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	select {
	case q.front <- n:
	default:
		if cap(q.front) == q.maxSize {
			// Drop the item; queue is full
			return ErrOverflow{}
		}
		newSize := 2 * cap(q.front)
		if newSize > q.maxSize {
			newSize = q.maxSize
		}
		oldCh := q.front
		close(oldCh)
		ch := make(chan int64, newSize)
		q.front = ch
		for v := range oldCh {
			q.front <- v
		}
		q.front <- n
	}
	return nil
}

// Dequeue returns the next item in the queue, or an error if the queue is empty.
func (q *Int64Queue) Dequeue() (int64, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	select {
	case v := <-q.front:
		return v, nil
	default:
		return 0, ErrUnderflow{}
	}
}
