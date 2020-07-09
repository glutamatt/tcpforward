package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	target string
	port   int
)

func init() {
	flag.StringVar(&target, "target", "", "the target (<host>:<port>)")
	flag.IntVar(&port, "port", 7757, "the tunnelthing port")
}

func main() {
	flag.Parse()
	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("could not start server on %d: %v", port, err)
	}
	fmt.Printf("server running on %d\n", port)
	for {
		client, err := incoming.Accept()
		if err != nil {
			log.Fatal("could not accept client connection", err)
		}
		go proxy(client)
	}
}

func proxy(client net.Conn) {
	fmt.Printf("client '%v' connected!\n", client.RemoteAddr())
	target, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatal("could not connect to target", err)
	}
	fmt.Printf("connection to server %v established!\n", target.RemoteAddr())

	go func() {
		io.Copy(target, client)
		target.Close()
	}()
	go func() {
		io.Copy(client, target)
		client.Close()
	}()
}
