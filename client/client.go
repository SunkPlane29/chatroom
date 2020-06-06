// THERE IS ROOM FOR IMPROVEMENT IN THIS FILE
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// LATER USE
var MSGHEADER = 16

type Client struct {
	Conn     net.Conn
	Username string
	Scanner  *bufio.Scanner
}

// Creates a new client given an existing connection
func newClient(c net.Conn) (*Client, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Type your username: ")
	scanner.Scan()
	username := scanner.Text()

	if username == "" || username == " " {
		return nil, errors.New("No valid username")
	}
	return &Client{
		Conn:     c,
		Username: username,
		Scanner:  scanner,
	}, nil
}

// Receives all messages and print them to the terminal
// ON PROGRESS: NEED TO RECEIVE USENAME HEADERS.
func (cl *Client) receiveMessages() {
	msgBuffer := make([]byte, MSGHEADER)
	for {
		var msg string
		_, err := cl.Conn.Read(msgBuffer)
		if err != nil {
			cl.Conn.Close()
			log.Fatal("FATAL ERROR, closing connection.")
		}
		msgUsername := strings.TrimSpace(string(msgBuffer))

		for {
			n, err := cl.Conn.Read(msgBuffer)
			if err != nil {
				cl.Conn.Close()
				log.Fatal("FATAL ERROR, closing connection.")
			}
			if n != len(msgBuffer) {
				msg += string(msgBuffer[:n])
				break
			}
			msg += string(msgBuffer)
		}

		fmt.Printf("\n<%s> %s", msgUsername, msg)
		fmt.Printf("\n<%s> ", cl.Username)
	}

}

// Scan lines on terminal to later sending that message
// ON PROGRESS.
func (cl *Client) sendMessages() {
	for {
		fmt.Printf("<%s> ", cl.Username)
		cl.Scanner.Scan()
		msg := cl.Scanner.Text()
		cl.Conn.Write([]byte(msg))

	}
}
