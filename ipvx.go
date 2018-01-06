// Package ipvx is a light-weight wrapper around the ipv4 and ipv6 packages, allowing
// applications to be more protocol agnostic.
//
// The Conn and PacketConn interfaces are largely drop-in replacements for the corresponding
// structs from the ipv6 package, and translate to ipv4 where needed.  Not all methods are
// implemented, where they seem to be tightly coupled to the specific IP version.
package ipvx

// A great deal more from ipv4 and ipv6 could be implemented here, requiring progressively
// more overhead in terms of translating struct fields between the v6 and v4 versions.

import "net"

type localAddrer interface {
	LocalAddr() net.Addr
}

func is4(c localAddrer) bool {
	if ipaddr, ok := c.LocalAddr().(*net.IPAddr); ok {
		if ipaddr.IP.To4() != nil {
			return true
		}
	}
	return false
}
