package ipvx_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/dnesting/ipvx"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func TestIs4(t *testing.T) {
	// This test goes a bit above and beyond.  Rather than creating test *net.IPAddr and similar
	// instances, we actually double-check that all of our assumptions about how the network stack
	// resolves addresses are true.
	type testcase struct {
		addr    string
		network string
		is4     bool
	}
	cases := []testcase{
		// This test only parses addresses, and so while these test cases are useful to check for
		// real-world behavior, we don't attempt to actually bounce them off of a sockets
		// implementation to test them.  And even if we did, some of these would be system-dependent.
		// {"", "ip", false},
		// {"", "ip4", true},
		// {"", "ip6", true},
		// {":123", "udp", false},
		// {":", "udp4", true},
		// {":", "udp6", true},
		// {":123", "tcp", false},
		// {":", "tcp4", true},
		// {":", "tcp6", true},

		// These should all be IPv4
		{"0.0.0.0", "ip", true},
		{"0.0.0.0", "ip4", true},
		{"0.0.0.0/0", "ip+net", true},
		{"0.0.0.0:0", "udp", true},
		{"0.0.0.0:0", "udp4", true},
		{"0.0.0.0:0", "tcp", true},
		{"0.0.0.0:0", "tcp4", true},
		{"127.0.0.1", "ip", true},
		{"127.0.0.1", "ip4", true},
		{"127.0.0.1/8", "ip+net", true},
		{"127.0.0.1:0", "udp", true},
		{"127.0.0.1:0", "udp4", true},
		{"127.0.0.1:0", "tcp", true},
		{"127.0.0.1:0", "tcp4", true},
		{"::ffff:127.0.0.1", "ip", true},
		{"::ffff:127.0.0.1", "ip4", true},
		{"::ffff:127.0.0.1/102", "ip+net", true}, // NB: IPv6 length; a /8 would yield ::/8 which is considered IPv6.
		{"[::ffff:127.0.0.1]:123", "udp", true},
		{"[::ffff:127.0.0.1]:123", "udp4", true},
		{"[::ffff:127.0.0.1]:123", "tcp", true},
		{"[::ffff:127.0.0.1]:123", "tcp4", true},

		// These should all be IPv6
		{"::", "ip", false},
		{"::", "ip6", false},
		{"::%en0", "ip6", false},
		{"::/64", "ip+net", false},
		{"[::]:0", "udp", false},
		{"[::]:0", "udp6", false},
		{"[::]:0", "tcp", false},
		{"[::]:0", "tcp6", false},
		{"::1", "ip", false},
		{"::1", "ip6", false},
		{"::1%en0", "ip6", false},
		{"::1/8", "ip+net", false},
		{"[::1]:0", "udp", false},
		{"[::1]:0", "udp6", false},
		{"[::1]:0", "tcp", false},
		{"[::1]:0", "tcp6", false},
	}

	for _, c := range cases {
		var addr net.Addr
		var err error

		switch c.network {
		case "ip", "ip4", "ip6":
			addr, err = net.ResolveIPAddr(c.network, c.addr)
		case "udp", "udp4", "udp6":
			addr, err = net.ResolveUDPAddr(c.network, c.addr)
		case "tcp", "tcp4", "tcp6":
			addr, err = net.ResolveTCPAddr(c.network, c.addr)
		case "ip+net":
			_, addr, err = net.ParseCIDR(c.addr)
		default:
			panic(fmt.Sprintf("unexpected network type %q in test", c.network))
		}

		if err != nil {
			t.Errorf("case %q %q: got unexpected error %s", c.network, c.addr, err)
			continue
		}
		is4 := ipvx.Is4(addr)
		if is4 != c.is4 {
			t.Errorf("case %q %q: expected is4=%v, got %v (ip=%v)", c.network, c.addr, c.is4, is4, ipvx.GetIP(addr))
		}
	}
}

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
