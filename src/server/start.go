package server

import (
	"dns-focus/src/common"
	"dns-focus/src/config"
	"dns-focus/src/resolver"
	"fmt"
	"log"
	"net"
)


func Start(dnsConfig *config.DnsConfig, serverMode common.ServerMode) {
	p, err := net.ListenPacket("udp", ":53")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	for {
		buf := make([]byte, 512)
		n, addr, err := p.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Connection error [%s]: %s\n", addr.String(), err)
			continue
		}
		go resolver.HandlePacket(p, addr, buf[:n], dnsConfig, serverMode)
	}
}