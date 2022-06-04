package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// A map to store current online users
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// A channel to broadcast messages
	Message chan string
}

// Create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
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

	// When starting, start the goroutine to monitor Message channel
	go this.ListenMessage()

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

// Func to monitor current Message channel, once there are any messages in the Channel, send it to all online users
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Connection create success.")

	user := NewUser(conn, this)

	// Broadcast user online message
	user.Online()

	// A channel to monitor whether the user is alive
	isLive := make(chan bool)

	// Recieve the messages sent by user and broadcast to all users
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err: ", err)
				return
			}

			// Extract the message from user (remove "\n")
			msg := string(buffer[:n-1])

			// Broadcast the message
			user.DoMessage(msg)

			// Any message sent by user means the user is alive
			isLive <- true
		}
	}()

	// Keep current goroutine alive
	for {
		select {
		case <-isLive:
			// The user is alive, should reset timer
			// Do nothing, to activate the select, update the timer
		case <-time.After(time.Second * 10):
			// Timeout, should force close the user
			user.sendMsg("You are timeout!")
			// isLive <- false
			close(user.C)
			// user.Offline()
			conn.Close()
			return
		}
	}

}

// Func to broadcast the user online message to Message Channel
func (this *Server) Broadcast(user *User, msg string) {
	sendMessage := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMessage
}
