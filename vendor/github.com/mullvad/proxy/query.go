package proxy

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/likexian/doh-go"
	"github.com/likexian/doh-go/dns"
	"github.com/mullvad/ipv6md"

	"github.com/mullvad/proxy/plain"
	"github.com/mullvad/proxy/xor"
	"github.com/mullvad/proxy/xorv2"
)

func Query(domains []string, verbose bool) []Proxy {
	var proxies []Proxy

	for _, d := range domains {
		if verbose {
			log.Printf("Querying AAAA record on %s for target addresses\n", d)
		}

		p, err := queryDomain(d)
		if err != nil {
			log.Printf("Failed to query %s for target addresses, %v\n", d, err)
			continue
		}

		proxies = append(proxies, p...)
	}

	return proxies
}

func queryDomain(domain string) ([]Proxy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := doh.Use(doh.CloudflareProvider, doh.GoogleProvider)

	r, err := c.Query(ctx, dns.Domain(domain), dns.TypeAAAA)
	if err != nil {
		return nil, err
	}
	c.Close()

	var proxies []Proxy
	for _, a := range r.Answer {
		ip := net.ParseIP(a.Data)
		if ip == nil {
			log.Printf("Unable to parse IP address from %s\n", a.Data)
			continue
		}

		typ, err := ipv6md.GetType(ip)
		if err != nil {
			log.Printf("Unable to determine type, %v\n", err)
			continue
		}

		switch typ {
		case ipv6md.AddrPort:
			p, err := plain.New(ip)
			if err != nil {
				log.Printf("Unable to decode address and port, %v\n", err)
				continue
			}
			proxies = append(proxies, p)
		case ipv6md.AddrPortXOR:
			p, err := xor.New(ip)
			if err != nil {
				log.Printf("Unable to decode address and port with XOR, %v\n", err)
				continue
			}
			proxies = append(proxies, p)
		case ipv6md.AddrPortXORV2:
			p, err := xorv2.New(ip)
			if err != nil {
				log.Printf("Unable to decode address and port with XOR v2, %v\n", err)
				continue
			}
			proxies = append(proxies, p)
		}
	}

	return proxies, nil
}
