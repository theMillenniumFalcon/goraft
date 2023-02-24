package main

import "github.com/theMillenniumFalcon/goraft/cache"

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
