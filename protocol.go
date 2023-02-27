package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	Cmd        Command
	Key        []byte
	Value      []byte
	timeToLive time.Duration
}

func (m *Message) ToBytes() []byte {
	switch m.Cmd {

	case CMDSet:
		cmd := fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.timeToLive)
		return []byte(cmd)

	case CMDGet:
		cmd := fmt.Sprintf("%s %s", m.Cmd, m.Key)
		return []byte(cmd)

	default:
		panic("unknown command")
	}
}

func parseMessage(raw []byte) (*Message, error) {
	var (
		rawString = string(raw)
		parts     = strings.Split(rawString, " ")
	)

	if len(parts) < 2 {
		/// respond
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			/// respond
			return nil, errors.New("invalid SET command")
		}
		msg.Value = []byte(parts[2])
		timeToLive, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid SET timeToLive")
		}
		msg.timeToLive = time.Duration(timeToLive)
	}

	return msg, nil
}
