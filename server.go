package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

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
	var (
		rawString = string(rawCmd)
		parts     = strings.Split(rawString, " ")
	)

	if len(parts) == 0 {
		/// respond
		log.Println("invalid command")
		return
	}

	cmd := Command(parts[0])
	if cmd == CMDSET {
		if len(parts) != 4 {
			// respond
			log.Println("invalid SET command")
			return
		}

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println("invalid SET command")
			return
		}

		msg := MSGSet{
			Key:   []byte(parts[1]),
			Value: []byte(parts[2]),
			TTL:   time.Duration(ttl),
		}

		if err := s.handleSetCommand(conn, msg); err != nil {
			// respond
			return
		}
	}
}

func (s *Server) handleSetCommand(conn net.Conn, msg MSGSet) error {
	fmt.Println("Handling the set command: ", msg)

	return nil
}
