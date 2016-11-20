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

	msg, err := osc.NewMessage(*address)
	if err != nil {
		log.Fatal(err)
	}
	if err := addArgs(msg, flag.Args()); err != nil {
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
		if err := msg.WriteString(data); err != nil {
			return errors.Wrapf(err, "could not add %s to message", data)
		}
	case "i":
		i, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "could not parse integer from %s", data)
		}
		if err := msg.WriteInt32(int32(i)); err != nil {
			return errors.Wrap(err, "could not add integer to message")
		}
	}
	return nil
}
