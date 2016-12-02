package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/scgolang/osc"
)

func main() {
	var (
		host    = flag.String("h", "127.0.0.1", "host")
		port    = flag.Int("p", 57120, "port")
		address = flag.String("a", "/foo", "OSC address")
	)
	flag.Parse()

	msg := osc.Message{Address: *address}
	if err := addArgs(&msg, flag.Args()); err != nil {
		log.Fatal(err)
	}

	// Send the message.
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*host, strconv.Itoa(*port)))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := osc.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Send(msg); err != nil {
		log.Fatal(err)
	}
}

// addArgs adds arguments to an osc message.
func addArgs(msg *osc.Message, args []string) error {
	for _, arg := range args {
		if err := addArg(msg, arg); err != nil {
			return errors.Wrapf(err, "could not add arg '%s' to msg", arg)
		}
	}
	return nil
}

const argsep = ":"

// addArg adds a single argument to an osc message.
// It expects osc message arguments to be formatted as <typetag>:<data>
func addArg(msg *osc.Message, arg string) error {
	pieces := strings.Split(arg, argsep)
	typetag, data := pieces[0], pieces[1]
	switch typetag {
	default:
		return errors.New("unsupported typetag: " + typetag)
	case "s":
		msg.Arguments = append(msg.Arguments, osc.String(data))
	case "i":
		i, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "could not parse integer from %s", data)
		}
		msg.Arguments = append(msg.Arguments, osc.Int(i))
	}
	return nil
}
