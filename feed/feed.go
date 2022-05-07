package feed

import "sync"

type message[T any] struct {
	data  T
	ready chan struct{}
	final bool
	next  *message[T]
}

type Feed[T any] struct {
	mutex sync.Mutex
	head  *message[T]
}

func NewFeed[T any]() *Feed[T] {
	return &Feed[T]{
		head: &message[T]{
			ready: make(chan struct{}),
			next:  nil,
			final: false,
		},
	}
}

func (f *Feed[T]) Publish(data T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.head.data = data
	f.head.next = &message[T]{
		ready: make(chan struct{}),
	}
	close(f.head.ready)
	f.head = f.head.next
}

func (f *Feed[T]) Finish(data T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.head.data = data
	f.head.next = &message[T]{
		ready: make(chan struct{}),
		final: true,
	}
	close(f.head.next.ready)
	close(f.head.ready)
}

func (f *Feed[T]) Next() *message[T] {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.head
}

func (f *Feed[T]) Wait() chan struct{} {
	return f.head.ready
}
