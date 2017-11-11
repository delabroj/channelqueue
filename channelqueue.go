// Package channelqueue implements a channel-based queue that grows as needed up
// to a maximum length.
package channelqueue

// ErrOverflow indicates that a value was enqueued to a queue that had already
// been filled to its max size.
type ErrOverflow struct{}

func (e ErrOverflow) Error() string {
	return "max queue size exceeded, value dropped"
}

// ErrUnderflow indicates that Dequeue was called on an empty queue.
type ErrUnderflow struct{}

func (e ErrUnderflow) Error() string {
	return "queue empty"
}
