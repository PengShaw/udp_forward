package forward

import (
	"errors"
	"net"
	"strings"

	"github.com/PengShaw/GoUtilsKit/logger"
	"github.com/PengShaw/GoUtilsKit/socket"
)

func checkAddrInput(i string) (string, string, error) {
	sp := strings.SplitN(i, ":", 2)
	if len(sp) != 2 {
		return "", "", errors.New("string format should be like: udp:0.0.0.0:514 or tcp:0.0.0.0:514 or unix:/path/to/unix.sock")
	}

	var err error
	switch sp[0] {
	case "udp":
		_, err = net.ResolveUDPAddr("udp", sp[1])
	case "tcp":
		_, err = net.ResolveTCPAddr("tcp", sp[1])
	case "unix":
		_, err = net.ResolveUnixAddr("unix", sp[1])
	default:
		err = errors.New("supported protocol should be: udp or tcp or unix")
	}

	return sp[0], sp[1], err
}

func listen(network, address string, mtu int, ch chan<- []byte) {
	switch network {
	case "udp":
		socket.RunUDPServer(address, mtu, ch)
	case "tcp":
		socket.RunTCPServer(address, mtu, ch)
	case "unix":
		socket.RunUnixServer(address, mtu, ch)
	}
}

func Run(s string, ds []string, v bool, vv bool, mtu int) {
	logger.SetLevel(logger.LevelError)
	if v {
		logger.SetLevel(logger.LevelInfo)
	}
	if vv {
		logger.SetLevel(logger.LevelDebug)
	}

	network, addr, err := checkAddrInput(s)
	chSource := make(chan []byte)
	if err != nil {
		logger.Errorf("Check listen input failed: %s", err)
		return
	}
	go listen(network, addr, mtu, chSource)

	chDestinations := []chan []byte{}
	for _, d := range ds {
		network, addr, err := checkAddrInput(d)
		chDestination := make(chan []byte)
		if err != nil {
			logger.Errorf("Check destination input failed: %s", err)
			return
		}
		go socket.RunSocketClient(network, addr, chDestination)
		chDestinations = append(chDestinations, chDestination)
	}

	for {
		data := <-chSource
		for _, ch := range chDestinations {
			ch <- data
		}
	}
}
