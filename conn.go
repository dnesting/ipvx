package ipvx

import (
	"net"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Conn represents the set of ipv{4,6}.Conn methods that are protocol-agnostic.
// See golang.org/x/net/ipv6.Conn for details about each of these methods.
type Conn interface {
	HopLimit() (int, error)
	SetHopLimit(hoplim int) error

	// To4 returns the underlying ipv4.Conn, or nil if this is an IPv6 connection.
	To4() *ipv4.Conn
	// To6 returns the underlying ipv6.Conn, or nil if this is an IPv4 connection.
	To6() *ipv6.Conn
}

type conn6 ipv6.Conn

func (c *conn6) To4() *ipv4.Conn { return nil }
func (c *conn6) To6() *ipv6.Conn { return (*ipv6.Conn)(c) }

type conn4 ipv4.Conn

func (c *conn4) To4() *ipv4.Conn              { return (*ipv4.Conn)(c) }
func (c *conn4) To6() *ipv6.Conn              { return nil }
func (c *conn4) HopLimit() (int, error)       { return c.To4().TTL() }
func (c *conn4) SetHopLimit(hoplim int) error { return c.To4().SetTTL(hoplim) }

// NewConn creates an IP protocol-agnostic Conn instance from c.
// If c does not seem to be an IPv4 connection, it is assumed to be IPv6.
func NewConn(c net.Conn) Conn {
	if is4(c) {
		return (*conn4)(ipv4.NewConn(c))
	}
	return (*conn6)(ipv6.NewConn(c))
}
