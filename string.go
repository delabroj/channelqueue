package channelqueue

import "sync"

// StringQueue holds the current channel and the maximum length of the channel.
type StringQueue struct {
	front   chan string
	maxSize int
	mux     sync.Mutex
}

// MakeStringQueue returns a pointer to an initialized queue struct.
func MakeStringQueue(n int) *StringQueue {
	ch := make(chan string, 1)
	return &StringQueue{front: ch, maxSize: n}
}

// Enqueue adds the value to the queue's channel. If the channel is full, the
// channel is replaced with another channel twice as large until maxSize is reached.
func (q *StringQueue) Enqueue(s string) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	select {
	case q.front <- s:
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
		ch := make(chan string, newSize)
		q.front = ch
		for v := range oldCh {
			q.front <- v
		}
		q.front <- s
	}
	return nil
}

// Dequeue returns the next item in the queue, or an error if the queue is empty.
func (q *StringQueue) Dequeue() (string, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	select {
	case v := <-q.front:
		return v, nil
	default:
		return "", ErrUnderflow{}
	}
}
