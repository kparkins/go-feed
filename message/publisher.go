package message

import (
	"errors"
	"sync"
)

type message[T any] struct {
	data     T
	ready    chan struct{}
	finished bool
	next     *message[T]
}

type Publisher[T any] struct {
	mutex *sync.Mutex
	head  *message[T]
}

func NewPublisher[T any]() *Publisher[T] {
	m := &message[T]{
		ready:    make(chan struct{}),
		next:     nil,
		finished: false,
	}
	return &Publisher[T]{
		mutex: &sync.Mutex{},
		head:  m,
	}
}

func (f *Publisher[T]) Publish(data T) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if f.head.finished {
		return errors.New("tried to publish to finished feed")
	}
	f.head.data = data
	f.head.next = &message[T]{
		ready:    make(chan struct{}),
		finished: false,
	}
	close(f.head.ready)
	f.head = f.head.next
	return nil
}

func (f *Publisher[T]) Finish() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.head.finished = true
	close(f.head.ready)
}

func (f *Publisher[T]) Subscribe() *Feed[T] {
	return NewFeed(f)
}
