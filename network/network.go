// Package network defines types to represent an IPX network.
package network

import (
	"github.com/fragglet/ipxbox/ipx"
)

// Network represents the concept of an IPX network.
type Network interface {
	// NewNode creates a new network node.
	NewNode() Node
}

// Node represents a node attached to an IPX network.
type Node interface {
	ipx.ReadWriteCloser

	// Address returns the IPX address of the node.
	Address() ipx.Addr

	// GetProperty populates the given value based on its type. Since
	// network implementations may consist of many layers, this will
	// query through the layers to fetch the property. If successful,
	// true is returned.
	GetProperty(value interface{}) bool
}

func NodeAddress(n Node) ipx.Addr {
	var result ipx.Addr
	if !n.GetProperty(&result) {
		return ipx.AddrNull
	}
	return result
}
