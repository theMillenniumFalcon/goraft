package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Status byte

type Command byte

type CommandGet struct {
	Key []byte
}

type CommandSet struct {
	Key        []byte
	Value      []byte
	TimeToLive int
}

type ResponseSet struct {
	Status Status
}

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)

	binary.Write(buf, binary.LittleEndian, int64(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int64(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, int64(c.TimeToLive))

	return buf.Bytes()
}

func ParseCommand(r io.Reader) {
	var cmd Command
	binary.Read(r, binary.LittleEndian, &cmd)

	switch cmd {
	case CmdSet:
		set := ParseSetCommand(r)
		fmt.Println(set)
	}
}

func ParseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLength int64
	binary.Read(r, binary.LittleEndian, keyLength)
	cmd.Key = make([]byte, keyLength)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	return cmd
}
