package channelqueue

import "testing"

func TestMakeQueue(t *testing.T) {
	queue := MakeQueue(3)
	if cap(queue.front) != 1 {
		t.Error("cap = ", cap(queue.front))
	}
	if len(queue.front) != 0 {
		t.Error("len = ", len(queue.front))
	}
}

func TestEnqueue_Overflow(t *testing.T) {
	queue := MakeQueue(2)
	queue.Enqueue(1)
	queue.Enqueue(1)
	err := queue.Enqueue(1)
	if _, ok := err.(ErrOverflow); !ok {
		t.Error("max queue length exceeded without overflow error")
	}
}

func TestEnqueue_MaxSize(t *testing.T) {
	queue := MakeQueue(3)
	queue.Enqueue(0)
	queue.Enqueue(0)
	queue.Enqueue(0)
	if cap(queue.front) != 3 {
		t.Error("cap = ", cap(queue.front))
	}
}

func TestDequeue(t *testing.T) {
	queue := MakeQueue(2)
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	x, err := queue.Dequeue()
	if x != 1 || err != nil {
		t.Fail()
	}
	x, err = queue.Dequeue()
	if x != 2 || err != nil {
		t.Fail()
	}
	x, err = queue.Dequeue()
	if x == 3 || err == nil {
		t.Error("max queue length exceeded")
	}
}

func TestDequeue_Underflow(t *testing.T) {
	queue := MakeQueue(2)
	_, err := queue.Dequeue()
	if _, ok := err.(ErrUnderflow); !ok {
		t.Error("underflow error not indicated")
	}
}

func TestEnqueueDequeueCycle(t *testing.T) {
	queue := MakeQueue(2)
	ch := make(chan int, queue.maxSize)
	for j := 0; j < 2; j++ {
		for i := 0; i < 3; i++ {
			err := queue.Enqueue(i)
			if err == nil {
				ch <- i
			}
		}
		for i := 0; i < 3; i++ {
			v, err := queue.Dequeue()
			if err == nil {
				if v != <-ch {
					t.Fail()
				}
			}
		}
	}
}
