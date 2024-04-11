package forward

import (
	"fmt"
	"net"
	"net/netip"
)

type Destination struct {
	ch   chan []byte
	addr netip.AddrPort
}

func Run(s netip.AddrPort, ds []netip.AddrPort, t bool) error {
	chErr := make(chan error)
	defer close(chErr)
	go runErrPrint(chErr)

	destinations := []Destination{}
	for _, d := range ds {
		destination := Destination{
			ch:   make(chan []byte),
			addr: d,
		}
		defer close(destination.ch)
		go runDestination(destination, chErr)
		destinations = append(destinations, destination)
	}

	return runSource(s, destinations, chErr, t)
}

func runErrPrint(ch <-chan error) {
	for {
		err := <-ch
		fmt.Printf("[Error] %s \n", err)
	}
}

func runDestination(d Destination, chErr chan<- error) {
	addr := net.UDPAddrFromAddrPort(d.addr)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		chErr <- fmt.Errorf("Connect destination(%s) failed: %s", addr.String(), err)
		return
	}
	defer conn.Close()

	for {
		data := <-d.ch
		_, err = conn.Write(data)
		if err != nil {
			chErr <- fmt.Errorf("Write data(%s) to destination(%s) failed: %s", data, addr.String(), err)
		}
	}
}

func runSource(s netip.AddrPort, ds []Destination, chErr chan<- error, t bool) error {
	addr := net.UDPAddrFromAddrPort(s)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Printf("Listen: <%s> \n", addr.String())

	data := make([]byte, 4096)
	for {
		n, err := conn.Read(data)
		if err != nil {
			chErr <- fmt.Errorf("Received data failed: %s", err)
		}
		for _, d := range ds {
			d.ch <- data[:n]
		}
		if t {
			fmt.Printf("Received data: %s \n", data[:n])
		}
	}
}
