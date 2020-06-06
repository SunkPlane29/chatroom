// For package readability, it's better to later change some things to show the least stuff possible
//  on main.go file
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	address = "127.0.0.1:8000"
	wg      = sync.WaitGroup{}
)

func main() {
	c, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("FATAL ERROR on main: %s", err.Error())
	}
	client, err := newClient(c)
	if err != nil {
		fmt.Println(err.Error())
		client, err = newClient(c)
		if err != nil {
			panic(err)
		}
	}
	c.Write([]byte(client.Username))

	wg.Add(2)
	go client.sendMessages()
	go client.receiveMessages()
	wg.Wait()
}
