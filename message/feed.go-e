package message

import "errors"

type Feed[T any] struct {
	message *message[T]
}

func NewFeed[T any](feed *Publisher[T]) *Feed[T] {
	return &Feed[T]{
		message: feed.head,
	}
}

func (s *Feed[T]) Value() T {
	return s.message.data
}

func (s *Feed[T]) Wait() chan struct{} {
	return s.message.ready
}

func (s *Feed[T]) Next() error {
	if s.message.finished {
		return errors.New("called Next() on a Finished() feed")
	}
	s.message = s.message.next
	return nil
}

func (s *Feed[T]) HasNext() bool {
	return !s.message.finished
}
