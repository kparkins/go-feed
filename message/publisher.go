package message

import (
	"sync"
)

type message[T any] struct {
	data     T
	ready    chan struct{}
	finished bool
	next     *message[T]
}

type Publisher[T any] struct {
	mutex sync.Mutex
	head  *message[T]
}

func NewPublisher[T any]() *Publisher[T] {
	m := &message[T]{
		ready:    make(chan struct{}),
		next:     nil,
		finished: false,
	}
	return &Publisher[T]{
		head: m,
	}
}

func (f *Publisher[T]) Publish(data T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.head.data = data
	f.head.next = &message[T]{
		ready: make(chan struct{}),
	}
	close(f.head.ready)
	f.head = f.head.next
}

func (f *Publisher[T]) Finish(data T) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.head.data = data
	f.head.next = &message[T]{
		ready:    make(chan struct{}),
		finished: true,
	}
	close(f.head.next.ready)
	close(f.head.ready)
}

func (f *Publisher[T]) Subscribe() *Feed[T] {
	return NewFeed(f)
}
