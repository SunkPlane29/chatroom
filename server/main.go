// For package readability, it's better to later change some things to show the least stuff possible
//  on main.go file
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	myServer := newServer()
	listener, err := myServer.startServer()
	if err != nil {
		log.Fatalf("FATAL ERROR on main: %s", err)
	}
	fmt.Printf("Listening on %s\n", myServer.Port)

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Changing the newClient function to just passing c as argument to make it less redundent.
		go func(c net.Conn) {
			myServer.Clients[c] = newClient(c)
		}(c)
		fmt.Printf("Connection accepted from %s\n", c.RemoteAddr())

		go myServer.receiveMessages(c)

		// Remover o cliente do map de clientes e talvez achar um identificador melhor para ele.
	}
}
