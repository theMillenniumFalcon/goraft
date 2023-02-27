package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/theMillenniumFalcon/goraft/cache"
	"github.com/theMillenniumFalcon/goraft/client"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":4000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	options := ServerOptions{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2)
		client, err := client.New(":4000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			SendCommand(client)
		}
		client.Close()
		time.Sleep(time.Second * 1)
	}()

	server := NewServer(options, cache.New())
	server.Start()
}

func SendCommand(c *client.Client) {
	_, err := c.Set(context.Background(), []byte("apple"), []byte("banana"), 0)
	if err != nil {
		log.Fatal(err)
	}
}
