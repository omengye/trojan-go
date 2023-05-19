package socks

import "github.com/p4gefau1t/trojan-go/config"

type Config struct {
	LocalHost  string `json:"local_addr" yaml:"local-addr"`
	LocalPort  int    `json:"local_port" yaml:"local-port"`
	UDPTimeout int    `json:"udp_timeout" yaml:"udp-timeout"`

	RemoteHost     string `json:"remote_addr" yaml:"remote-addr"`
	RemotePort     int    `json:"remote_port" yaml:"remote-port"`
	RemoteUsername string `json:"remote_username" yaml:"remote-username"`
	RemotePassword string `json:"remote_password" yaml:"remote-password"`
}

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return &Config{
			UDPTimeout: 60,
		}
	})
}
