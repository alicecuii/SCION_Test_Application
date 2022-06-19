package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"inet.af/netaddr"
	"os"
	"time"
)

func main() {
	var listen pan.IPPortValue
	var err error
	flag.Var(&listen, "listen", "[Server] local IP:port to listen on")
	remoteAddr := flag.String("remote", "", "[Client] Remote (i.e. the server's) SCION Address (e.g. 17-ffaa:1:1,[127.0.0.1]:12345)")
	flag.Parse()
	if (listen.Get().Port() > 0) == (len(*remoteAddr) > 0) {
		check(fmt.Errorf("either specify -listen for server or -remote for client"))
	}

	if listen.Get().Port() > 0 {
		err = runServer(listen.Get())
		check(err)
	} else {
		err = runClient(*remoteAddr)
		check(err)
	}
}
func runClient(address string) error {
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return err
	}
	conn, err := pan.DialUDP(context.Background(), netaddr.IPPort{}, addr, nil, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	message := "Hi Server, this is client speaking!"
	nBytes, err := conn.Write([]byte(fmt.Sprintf("%s", message)))
	if err != nil {
		return err
	}
	fmt.Printf("I am a Client, I wrote %d bytes. The message is: %s \n", nBytes, message)

	buffer := make([]byte, 16*1024)
	if err = conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return err
	}
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	data := buffer[:n]
	fmt.Printf("Received reply: %s\n", data)
	return nil
}
func runServer(listen netaddr.IPPort) error {
	conn, err := pan.ListenUDP(context.Background(), listen, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr())

	buffer := make([]byte, 16*1024)
	for {
		n, from, err := conn.ReadFrom(buffer)
		if err != nil {
			return err
		}
		data := buffer[:n]
		fmt.Printf("Received %s: %s\n", from, data)
		msg := fmt.Sprintf("This is server speaking, I have received your message!")
		n, err = conn.WriteTo([]byte(msg), from)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote %d bytes.\n", n)
	}
}

// Check just ensures the error is nil, or complains and quits
func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "Fatal error:", e)
		os.Exit(1)
	}
}
