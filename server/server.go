package main

import (
	"fmt"
	"net"
	"strings"
)

func makeHeader(header string, headerLength int) []byte {
	var h []byte
	for i := 0; i < len(header); i++ {
		h = append(h, byte(header[i]))
	}
	for j := 0; j < headerLength-len(header); j++ {
		var space byte = 32
		h = append(h, space)
	}
	return h
}

type Server struct {
	Port    string
	Address string
	Clients map[net.Conn]*Client // May want to set another thing as key
	Rooms   map[string]*Room
}

// Creates a server struct and returns a pointer to it
func newServer() *Server {
	return &Server{
		Port:    ":8000",
		Address: "127.0.0.1",
		Clients: make(map[net.Conn]*Client),
		Rooms:   make(map[string]*Room),
	}
}

// Starts the server
func (s *Server) startServer() (net.Listener, error) {
	listener, err := net.Listen("tcp", s.Address+s.Port)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (s *Server) removeClient(c net.Conn) {
	c.Close()
	delete(s.Clients, c)
}

// Get Client based on key which is the connection
func (s *Server) getClient(c net.Conn) *Client {
	return s.Clients[c]
}

// Some constants for room prefixes
const (
	CREATEPREF string = "/create"
	DELETEPREF string = "/delete"
	LISTPREF   string = "/list"
	JOINPREF   string = "/join"
	EXITPREF   string = "/exit"
)

// Receive the messages from a given client
func (s *Server) receiveMessages(c net.Conn) {
	buffer := make([]byte, 16)
	for {
		var msg string

		for {
			n, err := c.Read(buffer)
			if err != nil {
				s.removeClient(c)
				fmt.Printf("Connection closed from %s\n", c.RemoteAddr())
				return
			}
			if n != len(buffer) {
				msg += string(buffer[:n])
				break
			}
			msg += string(buffer)
		}
		s.handleCommands(msg, c)
	}
}

// Handles the messages if there is any command in them or
// calling the default case if there is not.
func (s *Server) handleCommands(msg string, c net.Conn) {
	msgCommands := strings.Split(msg, " ")
	switch msgCommands[0] {
	case JOINPREF:
		if len(msgCommands) == 1 {
			return
		}
		go s.joinRoom(msgCommands[1], c)
		go s.broadcastMessage("joined the room", c)
	case DELETEPREF:
		if len(msgCommands) == 1 {
			return
		}
		go s.deleteRoom(msgCommands[1], c)
	case LISTPREF:
		header := makeHeader("server", 16)
		roomsList := s.listRooms()
		rooms := strings.Join(roomsList, ", ")
		c.Write(header)
		c.Write([]byte(fmt.Sprintf("These are the current rooms availble: %s", rooms)))
	case CREATEPREF:
		if len(msgCommands) == 1 {
			return
		}
		go s.createRoom(msgCommands[1], c)
	case EXITPREF:
		s.exitRoom(c)
	default:
		go s.broadcastMessage(msg, c)
	}
}

// Broadcast messages from a given client
// IN PROGRESS
func (s *Server) broadcastMessage(msg string, c net.Conn) {
	clientRoom := s.getClient(c).CurrentRoom
	if clientRoom == nil {
		fmt.Println("User not currently in a room.")
		return
	}
	username := s.getClient(c).Username
	for k, _ := range clientRoom.Clients {
		if k != c {
			k.Write([]byte(username))
			k.Write([]byte(msg))
		}
	}
}
