package main

import (
	"fmt"
	"go-feed/feed"
	"sync"
)

func main() {
	f := feed.NewFeed[int]()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		sub := feed.NewSubcription(f)
		go func(thread int, s *feed.Subscription[int]) {
			for s.HasNext() {
				<-s.Signal()
				value := s.Value()
				fmt.Println(value)
				s.Next()
			}
			wg.Done()
		}(i, sub)
	}
	for i := 0; i < 10; i++ {
		f.Publish(i)
	}
	f.Finish(11)
	wg.Wait()
}
