package server

import (
	resolver "dns-server/resolver"
	"dns-server/server/dto"
	"fmt"
	"log"
	"net"
)


func Start(dnsConfig *dto.DnsConfig) {
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
		go resolver.HandlePacket(p, addr, buf[:n], dnsConfig)
	}
}