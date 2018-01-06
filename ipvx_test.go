package ipvx_test

import (
	"net"

	"github.com/dnesting/ipvx"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func Example() {
	var conn net.Conn

	// Where you might have done something like this in the past:
	if ipaddr, ok := conn.LocalAddr().(*net.IPAddr); ok {
		if ipaddr.IP.To4() == nil {
			ipv6.NewConn(conn).SetHopLimit(2)
		} else {
			ipv4.NewConn(conn).SetTTL(2)
		}
	}

	// Now you can just do:
	ipvx.NewConn(conn).SetHopLimit(2)
}
