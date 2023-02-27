package proto

import (
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

func (r *CommandSet) Bytes() []byte {

}

func ParseCommand(r io.Reader) (any, error) {

}
