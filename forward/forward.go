package forward

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

var appLog = NewLog()

type Addr struct {
	protocol string
	addr     string
}

func NewAddr(s string) (*Addr, error) {
	sp := strings.SplitN(s, ":", 2)
	if len(sp) != 2 {
		return nil, errors.New("string format should be like: udp:0.0.0.0:514 or unix:/path/to/unix.sock")
	}
	addr := Addr{
		protocol: sp[0],
		addr:     sp[1],
	}

	var err error
	switch addr.protocol {
	case "udp":
		_, err = net.ResolveUDPAddr("udp", addr.addr)
	case "unix":
		_, err = net.ResolveUnixAddr("unix", addr.addr)
	default:
		err = errors.New("supported protocol should be: udp or unix")
	}
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (a *Addr) Listen(ch chan<- []byte) {
	switch a.protocol {
	case "udp":
		a.udpListen(ch)
	case "unix":
		a.unixListen(ch)
	}
}

func (a *Addr) udpListen(ch chan<- []byte) {
	conn, err := net.ListenPacket(a.protocol, a.addr)
	if err != nil {
		appLog.Errorf("listen for %s:%s failed: %s", a.protocol, a.addr, err)
		return
	}
	defer conn.Close()
	appLog.Infof("listen: <%s>", conn.LocalAddr().String())

	buf := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			appLog.Errorf("listen data for %s:%s failed: %s", a.protocol, a.addr, err)
			continue
		}
		appLog.Infof("received data from %s", addr.String())
		appLog.Debugf("received data(%s) from %s", buf, addr.String())
		ch <- buf
	}
}

func (a *Addr) unixListen(ch chan<- []byte) {
	l, err := net.Listen(a.protocol, a.addr)
	if err != nil {
		appLog.Errorf("listen for %s:%s failed: %s", a.protocol, a.addr, err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		appLog.Infof("listen: <%s>", l.Addr().String())
		if err != nil {
			appLog.Errorf("connection for %s:%s failed: %s \n", a.protocol, a.addr, err)
			continue
		}

		go func(c net.Conn, ch chan<- []byte) {
			defer c.Close()
			for {
				buf := make([]byte, 1024)
				_, err := c.Read(buf)
				if err != nil && err != io.EOF {
					appLog.Errorf("listen data for %s:%s failed: %s", a.protocol, a.addr, err)
					break
				}
				if err == io.EOF {
					break
				}
				ch <- buf
				appLog.Infof("received data from %s", c.LocalAddr().String())
				appLog.Debugf("received data(%s) from %s", buf, c.LocalAddr().String())
			}
		}(conn, ch)
	}
}

func (a *Addr) Send(ch <-chan []byte) {
	conn, err := net.Dial(a.protocol, a.addr)
	if err != nil {
		appLog.Errorf("build connection to %s:%s failed: %s", a.protocol, a.addr, err)
		return
	}
	defer conn.Close()

	for {
		data := <-ch
		_, err := conn.Write(data)
		if err != nil {
			appLog.Errorf("send data to %s:%s failed: %s", a.protocol, a.addr, err)
			appLog.Debugf("send data(%s) to %s:%s failed: %s", data, a.protocol, a.addr, err)
		}
	}
}

func Run(s string, ds []string, v bool, vv bool) error {
	appLog.SetLevel(Error)
	if v {
		appLog.SetLevel(Info)
	}
	if vv {
		appLog.SetLevel(Debug)
	}

	source, err := NewAddr(s)
	chSource := make(chan []byte)
	if err != nil {
		return fmt.Errorf("listen %s", err)
	}
	go source.Listen(chSource)

	chDestinations := []chan []byte{}
	for _, d := range ds {
		destination, err := NewAddr(d)
		chDestination := make(chan []byte)
		if err != nil {
			return fmt.Errorf("destination %s", err)
		}
		go destination.Send(chDestination)
		chDestinations = append(chDestinations, chDestination)
	}

	for {
		data := <-chSource
		for _, ch := range chDestinations {
			ch <- data
		}
	}
}
