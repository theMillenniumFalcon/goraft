package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/theMillenniumFalcon/goraft/cache"
	"github.com/theMillenniumFalcon/goraft/client"
	"github.com/theMillenniumFalcon/goraft/proto"
)

type ServerOptions struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOptions
	members map[*client.Client]struct{}
	cache   cache.Cacher
}

func NewServer(options ServerOptions, c cache.Cacher) *Server {
	return &Server{
		ServerOptions: options,
		cache:         c,
		members:       make(map[*client.Client]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
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

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}

	log.Println("connected to leader:", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)

	s.handleConn(conn)

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	//fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}

	// fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	case *proto.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *proto.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}

	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	// log.Printf("GET %s", cmd.Key)

	resp := proto.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = proto.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TimeToLive)
			if err != nil {
				log.Println("forward to member error:", err)
			}
		}
	}()

	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TimeToLive)); err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = proto.StatusOK
	_, err := conn.Write(resp.Bytes())

	return err
}
