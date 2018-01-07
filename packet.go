package ipvx

import (
	"net"
	"time"

	"golang.org/x/net/bpf"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// PacketConn represents the set of ipv{4,6}.PacketConn methods that are protocol-agnostic.
// See golang.org/x/net/ipv6.PacketConn for details about each of these methods.
type PacketConn interface {
	Close() error
	ExcludeSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error
	HopLimit() (int, error)
	IncludeSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error
	JoinGroup(ifi *net.Interface, group net.Addr) error
	JoinSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error
	LeaveGroup(ifi *net.Interface, group net.Addr) error
	LeaveSourceSpecificGroup(ifi *net.Interface, group, source net.Addr) error
	MulticastHopLimit() (int, error)
	MulticastInterface() (*net.Interface, error)
	MulticastLoopback() (bool, error)
	SetBPF(filter []bpf.RawInstruction) error
	SetDeadline(t time.Time) error
	SetHopLimit(hoplim int) error
	SetMulticastHopLimit(hoplim int) error
	SetMulticastInterface(ifi *net.Interface) error
	SetMulticastLoopback(on bool) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

	// To4 returns the underlying ipv4.PacketConn, or nil if this is an IPv6 PacketConn.
	To4() *ipv4.PacketConn
	// To6 returns the underlying ipv6.PacketConn, or nil if this is an IPv4 PacketConn.
	To6() *ipv6.PacketConn
}

type pconn6 ipv6.PacketConn

func (c *pconn6) To4() *ipv4.PacketConn { return nil }
func (c *pconn6) To6() *ipv6.PacketConn { return (*ipv6.PacketConn)(c) }

type pconn4 ipv4.PacketConn

func (c *pconn4) To4() *ipv4.PacketConn                 { return (*ipv4.PacketConn)(c) }
func (c *pconn4) To6() *ipv6.PacketConn                 { return nil }
func (c *pconn4) HopLimit() (int, error)                { return c.To4().TTL() }
func (c *pconn4) SetHopLimit(hoplim int) error          { return c.To4().SetTTL(hoplim) }
func (c *pconn4) MulticastHopLimit() (int, error)       { return c.To4().MulticastTTL() }
func (c *pconn4) SetMulticastHopLimit(hoplim int) error { return c.To4().SetMulticastTTL(hoplim) }

// NewPacketConn creates an IP protocol-agnostic PacketConn instance from c.
// If c does not seem to be an IPv4 connection, it is assumed to be IPv6.
func NewPacketConn(c net.PacketConn) PacketConn {
	if Is4(c.LocalAddr()) {
		return (*pconn4)(ipv4.NewPacketConn(c))
	}
	return (*pconn6)(ipv6.NewPacketConn(c))
}
