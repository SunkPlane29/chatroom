package main

import (
	"fmt"
	"net"
	"strings"
)

// Room struct with a name and a client map
type Room struct {
	Name    string
	Creator net.Conn
	Clients map[net.Conn]*Client
}

// Create a room in the server given the roomName and it's
// creator. The last is for deleting them later.

// See if there is any room with the same name
func (s *Server) createRoom(roomName string, creator net.Conn) {
	fmt.Printf("User %s created room %s\n", strings.TrimSpace(s.getClient(creator).Username), roomName)
	s.Rooms[roomName] = &Room{
		Name:    roomName,
		Creator: creator,
		Clients: make(map[net.Conn]*Client),
	}
}

// Joining a client to a room, it should be a client
// sorta function, but for now a server suits well
func (s *Server) joinRoom(roomName string, c net.Conn) {
	r, ok := s.Rooms[roomName]
	if !ok {
		header := makeHeader("server", 16)
		roomsList := s.listRooms()
		rooms := strings.Join(roomsList, ", ")
		c.Write(header)
		c.Write([]byte(fmt.Sprintf("These are the current rooms availble: %s", rooms)))
		return
	}
	cl := s.getClient(c)
	cl.CurrentRoom = r
	r.Clients[c] = s.Clients[c]
	fmt.Printf("%s joined the rooom %s\n", strings.TrimSpace(s.getClient(c).Username), roomName)
}

// List the rooms in the server and returns a slice
// with their names
func (s *Server) listRooms() []string {
	var rooms []string
	for k, _ := range s.Rooms {
		rooms = append(rooms, k)
	}
	return rooms
}

// Later we must have a way to only be deletable by
// it's creator. And possibly later by an admin.
// May be good turning this bool into an error

// Making so every user in this room gets kicked out
// on delete
func (s *Server) deleteRoom(roomName string, c net.Conn) {
	if _, ok := s.Rooms[roomName]; !ok {
		header := makeHeader("server", 16)
		roomsList := s.listRooms()
		rooms := strings.Join(roomsList, ", ")
		c.Write(header)
		c.Write([]byte(fmt.Sprintf("These are the current rooms availble: %s", rooms)))
		return
	}

	if c != s.Rooms[roomName].Creator {
		header := makeHeader("server", 16)
		c.Write(header)
		c.Write([]byte("Only the creator can delete the room."))
		return
	}
	delete(s.Rooms, roomName)
}

// Exits a client from certain room, as joinRoom it maybe
// should be a client function but for now suits well
func (s *Server) exitRoom(c net.Conn) {
	cl := s.getClient(c)
	if cl.CurrentRoom == nil {
		fmt.Println("Not currently in a room")
		return
	}
	cl.CurrentRoom = nil
}
