package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
)

type Server struct {
	addr   string
	ls     net.Listener
	schema []byte
}

func New(addr string) (*Server, error) {
	ls, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}

	fmt.Println(ls.Addr())

	return &Server{
		ls:     ls,
		addr:   addr,
		schema: []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()-_=+[]{}|;:'\",.<>/?`~"),
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

		size, err := strconv.Atoi(string(reader[:n-2]))
		if err != nil {
			fmt.Println("Error:", err)
			conn.Write([]byte("wrong input\n"))
			return
		}

		conn.Write(s.Random(size))
	}
}

func (s *Server) Random(size int) []byte {
	message := make([]byte, size+1)

	for index := 0; index < len(message); index++ {
		message[index] = s.schema[rand.Intn(len(s.schema))]
	}

	message[len(message)-1] = '\n'

	return message
}

func (s *Server) Close() {
	s.ls.Close()
}

func main() {
	addr := ":5000"

	server, err := New(addr)
	if err != nil {
		log.Println(err)
	}

	defer server.Close()

	server.Listen()
}
