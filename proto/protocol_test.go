package proto

import (
	"bytes"
	"testing"
)

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:        []byte("Foo"),
		Value:      []byte("Bar"),
		TimeToLive: 2,
	}

	r := bytes.NewReader(cmd.Bytes())
}
