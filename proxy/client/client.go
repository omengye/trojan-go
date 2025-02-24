package client

import (
	"context"

	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/proxy"
	"github.com/p4gefau1t/trojan-go/tunnel/adapter"
	"github.com/p4gefau1t/trojan-go/tunnel/http"
	"github.com/p4gefau1t/trojan-go/tunnel/mux"
	"github.com/p4gefau1t/trojan-go/tunnel/router"
	"github.com/p4gefau1t/trojan-go/tunnel/simplesocks"
	"github.com/p4gefau1t/trojan-go/tunnel/socks"
	"github.com/p4gefau1t/trojan-go/tunnel/tls"
	"github.com/p4gefau1t/trojan-go/tunnel/transport"
	"github.com/p4gefau1t/trojan-go/tunnel/trojan"
	"github.com/p4gefau1t/trojan-go/tunnel/websocket"
)

const Name = "CLIENT"

// GenerateClientTree generate general outbound protocol stack
func GenerateClientTree(transportPlugin bool, muxEnabled bool, wsEnabled bool, routerEnabled bool) []string {
	clientStack := []string{transport.Name}
	if !transportPlugin {
		clientStack = append(clientStack, tls.Name)
	}
	if wsEnabled {
		clientStack = append(clientStack, websocket.Name)
	}
	clientStack = append(clientStack, trojan.Name)
	if muxEnabled {
		clientStack = append(clientStack, []string{mux.Name, simplesocks.Name}...)
	}
	if routerEnabled {
		clientStack = append(clientStack, router.Name)
	}
	return clientStack
}

func init() {
	proxy.RegisterProxyCreator(Name, func(ctx context.Context) (*proxy.Proxy, error) {
		cfg := config.FromContext(ctx, Name).(*Config)
		adapterServer, err := adapter.NewServer(ctx, nil)
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithCancel(ctx)

		root := &proxy.Node{
			Name:       adapter.Name,
			Next:       make(map[string]*proxy.Node),
			IsEndpoint: false,
			Context:    ctx,
			Server:     adapterServer,
		}

		root.BuildNext(http.Name).IsEndpoint = true
		root.BuildNext(socks.Name).IsEndpoint = true

		clientStack := GenerateClientTree(cfg.TransportPlugin.Enabled, cfg.Mux.Enabled, cfg.Websocket.Enabled, cfg.Router.Enabled)
		c, err := proxy.CreateClientStack(ctx, clientStack)
		if err != nil {
			cancel()
			return nil, err
		}
		s := proxy.FindAllEndpoints(root)
		return proxy.NewProxy(ctx, cancel, s, c), nil
	})
}
