package message

type Feed[T any] struct {
	message *message[T]
}

func NewFeed[T any](feed *Publisher[T]) *Feed[T] {
	return &Feed[T]{
		message: feed.head,
	}
}

func (s *Feed[T]) Value() T {
	var value T
	if s.message != nil {
		return s.message.data
	}
	return value
}

func (s *Feed[T]) Updated() chan struct{} {
	if s.message != nil {
		return s.message.ready
	}
	value := make(chan struct{})
	close(value)
	return value
}

func (s *Feed[T]) Next() bool {
	finished := s.message.finished
	if s.message != nil {
		s.message = s.message.next
	}
	return !finished
}

func (s *Feed[T]) Finished() bool {
	if s.message != nil {
		return s.message.finished
	}
	return true
}

func (s *Feed[T]) Unsubscribe() {
	s.message = nil
}
