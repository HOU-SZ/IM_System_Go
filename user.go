package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//Create a new user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

func (this *User) Online() {
	// Add the user to server's OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// Broadcast the user online message to server's Message Channel
	this.server.Broadcast(this, "online")
}

func (this *User) Offline() {
	// Delete the user from server's OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// Broadcast the user offline message to server's Message Channel
	this.server.Broadcast(this, "offline")
}

func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}
