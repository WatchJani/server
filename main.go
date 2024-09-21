package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	addr string
	ls   net.Listener
}

func New(addr string) (*Server, error) {
	ls, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}

	fmt.Println(ls.Addr())

	return &Server{
		ls:   ls,
		addr: addr,
	}, nil
}

func (s *Server) Listen() {
	for {
		conn, err := s.ls.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	reader := make([]byte, 4096)
	for {
		n, err := conn.Read(reader)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(string(reader[:n]))
	}
}

func (s *Server) Close() {
	s.ls.Close()
}

func main() {
	addr := "109.165.187.47:5000"

	addr = "localhost:"

	server, err := New(addr)
	if err != nil {
		log.Println(err)
	}

	defer server.Close()

	server.Listen()
}
