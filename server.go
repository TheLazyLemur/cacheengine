package main

import (
	"fmt"
	"log"
	"net"

	"github.com/thelazylemur/cacheengine/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader bool
}

type Server struct {
	ServerOpts ServerOpts
	cacher cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cacher: c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ServerOpts.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ServerOpts.ListenAddr)

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
		if err := conn.Close(); err != nil {
			log.Printf("close error: %s\n", err)
		}	
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error: %s\n", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		log.Println("error parsing command")
		return
	}

	if err := handleSetCommand(conn, msg); err != nil {
		log.Println("something went wrong while handling the SET command: ", msg)
		return
	}
}

func handleSetCommand(conn net.Conn, msg *Message) error {
	log.Println("handling the set command: ", msg)
	return nil
}
