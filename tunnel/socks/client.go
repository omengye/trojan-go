package socks

import (
	"context"

	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/tunnel/trojan"
	"golang.org/x/net/proxy"

	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/tunnel"
)

type Client struct {
	underlay  tunnel.Client
	proxyAddr *tunnel.Address
	auth      *proxy.Auth
}

func (c *Client) DialConn(addr *tunnel.Address, t tunnel.Tunnel) (tunnel.Conn, error) {
	dialer, err := proxy.SOCKS5("tcp", c.proxyAddr.String(), c.auth, proxy.Direct)
	if err != nil {
		return nil, common.NewError("failed to init socks dialer")
	}
	conn, err := dialer.Dial("tcp", addr.String())
	if err != nil {
		return nil, common.NewError("socks failed to dial using socks5 proxy").Base(err)
	}
	return &Conn{
		Conn: conn,
		metadata: &tunnel.Metadata{
			Command: Connect,
			Address: addr,
		},
	}, nil
}

func (c *Client) DialPacket(t tunnel.Tunnel) (tunnel.PacketConn, error) {
	conn, err := c.underlay.DialConn(nil, &Tunnel{})
	if err != nil {
		return nil, common.NewError("socks failed to dial using underlying tunnel").Base(err)
	}
	metadata := &tunnel.Metadata{
		Command: Associate,
		Address: &tunnel.Address{
			DomainName:  "UDP_CONN",
			AddressType: tunnel.DomainName,
		},
	}
	if err := metadata.WriteTo(conn); err != nil {
		return nil, common.NewError("socks failed to write udp associate").Base(err)
	}
	return &PacketConn{
		PacketConn: &trojan.PacketConn{
			Conn: conn,
		},
	}, nil
}

func (c *Client) Close() error {
	return c.underlay.Close()
}

func NewClient(ctx context.Context, underlay tunnel.Client) (*Client, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	socketServerAddr := tunnel.NewAddressFromHostPort("tcp", cfg.RemoteHost, cfg.RemotePort)
	auth := &proxy.Auth{
		User:     cfg.RemoteUsername,
		Password: cfg.RemotePassword,
	}
	log.Debug("socks client created")
	return &Client{
		underlay:  underlay,
		proxyAddr: socketServerAddr,
		auth:      auth,
	}, nil
}
