package net

import (
	"net"
	"strconv"
)

func GrabEphemeralPort() (port uint16, err error) {
	var listener net.Listener
	var portStr string
	var p int

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	defer listener.Close()

	_, portStr, err = net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return
	}

	p, err = strconv.Atoi(portStr)
	port = uint16(p)

	return
}
