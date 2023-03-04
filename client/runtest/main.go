package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/theMillenniumFalcon/goraft/client"
)

type Server struct {
	raft *raft.Raft
}

func main() {
	var (
		config             = raft.DefaultConfig()
		finiteStateMachine = &raft.MockFSM{}
		stableStore        = raft.NewInmemStore()
		logStore           = raft.NewInmemStore()
		snapShotStore      = raft.NewInmemSnapshotStore()
		timeout            = time.Second * 5
	)

	config.LocalID = "NPID"

	ips, err := net.LookupIP("localhost")
	if err != nil {
		log.Fatal(err)
	}
	if len(ips) == 0 {
		log.Fatalf("localhost did not resolve to any IPs")
	}
	addr := &net.TCPAddr{IP: ips[0], Port: 4000}

	fmt.Println(addr)

	tr, err := raft.NewTCPTransport(":3000", addr, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal("tcp net failed", err)
	}
	r, err := raft.NewRaft(config, finiteStateMachine, stableStore, logStore, snapShotStore, tr)
	if err != nil {
		log.Fatal("failed to create new raft: ", err)
	}

	fmt.Printf("%+v\n", r)

	select {}
}

func SendStuff() {
	c, err := client.New(":4000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		var (
			key   = []byte(fmt.Sprintf("key_%d", i))
			value = []byte(fmt.Sprintf("val_%d", i))
		)

		if err = c.Set(context.Background(), key, value, 0); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}
	c.Close()
}
