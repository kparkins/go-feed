package main

import (
	"context"
	"fmt"
	"go-message/message"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Publish(ctx context.Context, wg *sync.WaitGroup, pub *message.Publisher[int]) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			pub.Finish()
			goto Finished
		default:
			pub.Publish(i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			i++
		}
	}
Finished:
	fmt.Printf("reached %d messages sent", i)
	wg.Done()
}

func RunSubscriber(thread int, feed *message.Feed[int], wg *sync.WaitGroup) {
	for {
		select {
		case <-feed.Updated():
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			fmt.Println(feed.Value())
			if !feed.Next() {
				goto Done
			}
		}
	}
Done:
	wg.Done()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	pub := message.NewPublisher[int]()
	go func() {
		<-sigc
		cancel()
	}()
	fmt.Println("starting go-routines...")
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		stream := pub.Subscribe()
		go RunSubscriber(i, stream, &wg)
	}

	go Publish(ctx, &wg, pub)
	wg.Wait()
}
