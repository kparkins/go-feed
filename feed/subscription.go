package feed

type Subscription[T any] struct {
	message *message[T]
}

func NewSubscription[T any](feed *Feed[T]) *Subscription[T] {
	return &Subscription[T]{
		message: feed.head,
	}
}

func (s *Subscription[T]) Value() T {
	return s.message.data
}

func (s *Subscription[T]) Wait() chan struct{} {
	return s.message.ready
}

func (s *Subscription[T]) Next() {
	s.message = s.message.next
}

func (s *Subscription[T]) HasNext() bool {
	return !s.message.final
}
