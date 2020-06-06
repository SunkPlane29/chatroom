package main

import (
	"net"
)

type Client struct {
	Address     net.Addr
	Conn        net.Conn
	Username    string
	CurrentRoom *Room
}

// Create a client object and returns a pointer to it - FINDING A BETTER WAY TO DEAL WITH THIS USERNAME HEADER PROBLEM
func newClient(c net.Conn) *Client {
	// We can use the buffer as the header for usernames, every username is 16 long máx

	// Later finfing a way to make this better, it bad, but working at last

	// We are making a slice of bytes filled with spaces and then storing the username received to it
	// making an username with 16 characters with some spaces and storing that string to the client struct
	// it's not the best solution, it's bad.
	usernameBuffer := makeHeader(" ", 16)
	c.Read(usernameBuffer)
	username := string(usernameBuffer)

	return &Client{
		Address:     c.RemoteAddr(),
		Conn:        c,
		Username:    username,
		CurrentRoom: nil,
	}
}
