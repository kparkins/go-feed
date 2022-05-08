package main

import (
	"bufio"
	"context"
	"fmt"
	"go-message/message"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type CLIService struct {
	history []string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewCLIServce(ctx context.Context, cancel context.CancelFunc) *CLIService {
	return &CLIService{
		history: make([]string, 0),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *CLIService) readStdin(out chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		out <- input
	}
}

func (c *CLIService) handleInput(input string) {
	switch input {
	case "exit":
		c.cancel()
	case "help":
		fmt.Println("trying to help you")
	default:
	}
}

func (c *CLIService) Run() {
	input := make(chan string)
	go c.readStdin(input)

	for {
		select {
		case <-c.ctx.Done():
			return
		case line := <-input:
			c.handleInput(line)
		}

	}
}

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
			go func(thread int, feed *message.Feed[int]) {
				for feed.HasNext() {
					<-feed.Wait()
					value := feed.Value()
					fmt.Println(value)
					feed.Next()
				}
				wg.Done()
			}(i, stream)
		}
		for i := 0; i < 10; i++ {
			pub.Publish(i)
		}
		pub.Finish(11)

		wg.Done()
	}(&wg)
	cli := NewCLIServce(ctx, cancel)
	cli.Run()
	wg.Wait()
}
