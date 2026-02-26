package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var natsURL string
	var topic string
	var timeout time.Duration
	flag.StringVar(&natsURL, "nats-url", "nats://127.0.0.1:4222", "NATS URL")
	flag.StringVar(&topic, "topic", "episode.published.v1", "subject to subscribe")
	flag.DurationVar(&timeout, "timeout", 20*time.Second, "wait timeout")
	flag.Parse()

	conn, err := nats.Connect(natsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch := make(chan *nats.Msg, 1)
	sub, err := conn.Subscribe(topic, func(msg *nats.Msg) {
		select {
		case ch <- msg:
		default:
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "subscribe failed: %v\n", err)
		os.Exit(1)
	}
	defer sub.Unsubscribe()

	select {
	case msg := <-ch:
		fmt.Println(string(msg.Data))
		os.Exit(0)
	case <-time.After(timeout):
		fmt.Fprintf(os.Stderr, "timeout waiting for topic %s\n", topic)
		os.Exit(2)
	}
}
