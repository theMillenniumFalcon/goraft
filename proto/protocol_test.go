package proto

import (
	"fmt"
	"testing"
)

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:        []byte("Foo"),
		Value:      []byte("Bar"),
		TimeToLive: 2,
	}
	fmt.Println(cmd.Bytes())

	// r := bytes.NewReader(cmd.Bytes())
}
