package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

//Create a new user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// start the goroutine to monitor current channel
	go user.ListenMessage()

	return user
}

// func to monitor current user channel, once there are any message in the channel, send the message to client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
