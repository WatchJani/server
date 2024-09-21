package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
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
			conn.Close()
			break
		}

		size, err := strconv.Atoi(string(reader[:n-2]))
		if err != nil || size < 0 {
			fmt.Println("Error:", err)
			conn.Write([]byte("wrong input\n"))
			continue
		}

		log.Printf("%v | %v | %d\n", time.Now(), conn.RemoteAddr(), size)

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
	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	addr := ":5000"

	server, err := New(addr)
	if err != nil {
		log.Println(err)
	}

	defer server.Close()

	server.Listen()
}
