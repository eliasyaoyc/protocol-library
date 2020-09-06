package slice

import "sync"

// SliceQueue us an unbounded queue which uses a slice as underlying
type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

// NewSliceQueue returns an empty queue
func NewSliceQueue(n int) *SliceQueue {
	return &SliceQueue{
		data: make([]interface{}, n),
	}
}

func (s *SliceQueue) Enqueue(v interface{}) {
	s.mu.Lock()
	s.data = append(s.data, v)
	s.mu.Unlock()
}

func (s *SliceQueue) Dequeue() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) == 0 {
		return nil
	}
	i := s.data[0]
	s.data = s.data[1:]
	return i
}
