package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/scgolang/osc"
)

func main() {
	var (
		host    = flag.String("h", "127.0.0.1", "host")
		port    = flag.Int("p", 57120, "port")
		address = flag.String("a", "/foo", "OSC address")
	)
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, strconv.Itoa(*port)))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("addr %s\n", addr)
	conn, err := osc.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sending message to %s\n", *address)
	msg, err := osc.NewMessage(*address)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Send(msg); err != nil {
		log.Fatal(err)
	}
}
