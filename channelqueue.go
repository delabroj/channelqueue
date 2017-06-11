// Package channelqueue implements a channel-based queue that grows as needed up
// to a maximum length.
package channelqueue

import "sync"

// Queue holds the current channel and the maximum length of the channel.
type Queue struct {
	front   chan int
	maxSize int
	mux     sync.Mutex
}

// MakeQueue returns a pointer to an initialized queue struct.
func MakeQueue(n int) *Queue {
	ch := make(chan int, 1)
	return &Queue{front: ch, maxSize: n}
}

// ErrOverflow indicates that a value was enqueued to a queue that had already
// been filled to its max size.
type ErrOverflow struct{}

func (e ErrOverflow) Error() string {
	return "max queue size exceeded, value dropped"
}

// Enqueue adds the value to the queue's channel. If the channel is full, the
// channel is replaced with another channel twice as large until maxSize is reached.
func (q *Queue) Enqueue(n int) error {
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
		ch := make(chan int, newSize)
		q.front = ch
		for v := range oldCh {
			q.front <- v
		}
		q.front <- n
	}
	return nil
}

// ErrUnderflow indicates that Dequeue was called on an empty queue.
type ErrUnderflow struct{}

func (e ErrUnderflow) Error() string {
	return "queue empty"
}

// Dequeue returns the next item in the queue, or an error if the queue is empty.
func (q *Queue) Dequeue() (int, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	select {
	case v := <-q.front:
		return v, nil
	default:
		return 0, ErrUnderflow{}
	}
}
