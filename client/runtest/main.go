package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/theMillenniumFalcon/goraft/client"
)

func main() {
	SendStuff()
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
