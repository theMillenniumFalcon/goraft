package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/theMillenniumFalcon/goraft/cache"
	"github.com/theMillenniumFalcon/goraft/client"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":4000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	opts := ServerOptions{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	server := NewServer(opts, cache.New())
	server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			client, err := client.New(":4000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("val_%d", i))
			)

			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			fetchedValue, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(fetchedValue))

			client.Close()
		}(i)
	}
}
