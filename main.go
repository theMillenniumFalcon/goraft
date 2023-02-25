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

	// go func() {
	// 	time.Sleep(time.Second * 2) // sleep for 2 seconds
	// 	conn, err := net.Dial("tcp", ":4000")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	conn.Write([]byte("SET Foo Bar 250000000000"))

	// 	time.Sleep(time.Second * 2)
	// 	conn.Write([]byte("GET Foo"))

	// 	buf := make([]byte, 1000)
	// 	n, _ := conn.Read(buf)
	// 	fmt.Println(string(buf[:n]))
	// }()

	server := NewServer(options, cache.New())
	server.Start()
}
