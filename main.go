package main

import (
	"flag"

	"github.com/theMillenniumFalcon/goraft/cache"
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

	server := NewServer(options, cache.New())
	server.Start()
}
