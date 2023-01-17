package main

import (
	"context"
	"fmt"
	"go-message/message"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigc
		cancel()
	}()
	fmt.Println("starting go-routines...")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		pub := message.NewPublisher[int]()
		for i := 0; i < 10; i++ {
			wg.Add(1)
			stream := pub.Subscribe()
			go func(thread int, feed *message.Feed[int], wg *sync.WaitGroup) {
				for {
					select {
					case <-feed.Updated():
						fmt.Println(feed.Value())

						if !feed.Next() {
							goto Done
						}
					}
				}
			Done:
				wg.Done()
			}(i, stream, wg)
		}
		i := 0
		for {
			select {
			case <-ctx.Done():
				goto Finish
			default:
				pub.Publish(i)
				time.Sleep(time.Millisecond * 100)
				i++
			}
		}
	Finish:
		fmt.Printf("reached %d messages sent", i)
		pub.Finish()

		wg.Done()
	}(&wg)
	wg.Wait()
}
