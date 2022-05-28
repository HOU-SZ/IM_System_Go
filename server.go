package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// Create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Start() {
	// Socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	} else {
		fmt.Println("net.Listen success.")
	}

	// Finally close listener
	defer listener.Close()

	// Monitor connections
	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accpet error: ", err)
			continue
		}
		fmt.Println("Listener accept success.")
		go this.Handler(conn)

	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection create success.")
}
