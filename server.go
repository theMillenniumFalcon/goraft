package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/theMillenniumFalcon/goraft/cache"
)

type ServerOptions struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOptions

	cache cache.Cacher
}

func NewServer(options ServerOptions, _cache cache.Cacher) *Server {
	return &Server{
		ServerOptions: options,
		cache:         _cache,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("connn read error: %s\n", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse the command with error: ", err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCommand(conn, msg)
	case CMDGet:
		err = s.handleGetCommand(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command: ", err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg *Message) error {
	fmt.Println("Handling the set command: ", msg)

	if err := s.cache.Set(msg.Key, msg.Value, msg.timeToLive); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}

	_, err = conn.Write(val)

	return err
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
