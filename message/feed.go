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
	return s.message.data
}

func (s *Feed[T]) Updated() chan struct{} {
	return s.message.ready
}

func (s *Feed[T]) Next() bool {
	s.message = s.message.next
	return s.message.finished
}

func (s *Feed[T]) Finished() bool {
	return !s.message.finished
}

func (s *Feed[T]) Unsubscribe() {
	s.message = nil
}
