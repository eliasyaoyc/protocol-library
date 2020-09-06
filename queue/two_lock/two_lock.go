package two_lock

import "sync"

// TwoLock is a concurrent unbounded queue which uses two-lock concurrent queue algorithm
type TwoLock struct {
	head  *cnode
	tail  *cnode
	hlock sync.Mutex
	tlock sync.Mutex
}

type cnode struct {
	value interface{}
	next  *cnode
}

func NewTwoLock() *TwoLock {
	n := &cnode{}
	return &TwoLock{head: n, tail: n}
}

func (t *TwoLock) Enqueue(v interface{}) {
	n := &cnode{value: v}
	t.tlock.Lock()
	t.tail.next = n
	t.tail = n
	t.tlock.Unlock()
}

func (q *TwoLock) Dequeue() interface{} {
	q.hlock.Lock()
	n := q.head
	newHead := n.next
	if newHead == nil {
		q.hlock.Unlock()
		return nil
	}
	v := newHead.value
	newHead.value = nil
	q.head = newHead
	q.hlock.Unlock()
	return v
}
