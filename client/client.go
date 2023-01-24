// Package client implements a client for sending and receiving IPX frames
// from a server over UDP.
package client

import (
	"context"
	"errors"
	"net"

	"github.com/fragglet/ipxbox/ipx"
	"github.com/fragglet/ipxbox/network/pipe"
)

var (
	_ = (ipx.ReadWriteCloser)(&Client{})
)

// Client is an implementation of the ipx.ReadWriteCloser interface that
// sends and receives IPX frames to a server over UDP.
// This is *not* a complete implementation of the dosbox IPX protocol.
type Client struct {
	net.PacketConn
	ipx.ReadWriteCloser
	net.Addr
}

// Dial creates a new client for sending IPX frames to the server at the
// given address.
var Dial = DefaultDial

func DefaultDial(addr string) (*Client, error) {
	resolvedAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	PacketConn, err := net.DialUDP("udp4", nil, resolvedAddr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		PacketConn:      PacketConn,
		ReadWriteCloser: pipe.New(),
		Addr:            resolvedAddr,
	}
	go c.recvLoop()
	return c, nil
}

func (c *Client) recvLoop() {
	var buf [1500]byte
	defer c.ReadWriteCloser.Close()

	for {
		packetLen, _, err := c.PacketConn.ReadFrom(buf[:])
		if errors.Is(err, net.ErrClosed) {
			return
		} else if err != nil {
			// TODO: Log error?
			continue
		}

		p := &ipx.Packet{}
		if err := p.UnmarshalBinary(buf[:packetLen]); err != nil {
			// TODO: Log error?
			continue
		}
		if err := c.ReadWriteCloser.WritePacket(p); err != nil {
			// TODO: Log error?
		}
	}
}

func (c *Client) ReadPacket(ctx context.Context) (*ipx.Packet, error) {
	return c.ReadWriteCloser.ReadPacket(ctx)
}

func (c *Client) WritePacket(packet *ipx.Packet) error {
	packetBytes, err := packet.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = c.PacketConn.WriteTo(packetBytes, c.Addr)
	return err
}

func (c *Client) Close() error {
	c.ReadWriteCloser.Close()
	return c.PacketConn.Close()
}
