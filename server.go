package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/theMillenniumFalcon/goraft/cache"
	"github.com/theMillenniumFalcon/goraft/proto"
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

func NewServer(opts ServerOptions, c cache.Cacher) *Server {
	return &Server{
		ServerOptions: opts,
		cache:         c,
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
	defer conn.Close()

	fmt.Println("connection made: ", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error: ", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}

	fmt.Println("connection closed: ", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)

	resp := proto.ResponseSet()
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TimeToLive)); err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())

		return err
	}

	resp.Status = proto.StatusOK
}
