package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/mullvad/proxy"
	"github.com/mullvad/proxy/typ"
)

var appVersion string

func main() {
	listenAddr := flag.String("listen-addr", "127.0.0.1:1080", "listen address and port")
	sourceDomains := flag.String("source-domains", "frakta.eu", "source domains")
	verbose := flag.Bool("verbose", false, "verbose output")
	version := flag.Bool("version", false, "display version and exit")
	usePlainProxies := flag.Bool("use-plain-proxies", true, "use plain proxies")
	useXORProxies := flag.Bool("use-xor-proxies", true, "use xor proxies")
	useXORV2Proxies := flag.Bool("use-xor-v2-proxies", true, "use xor v2 proxies")
	flag.Parse()

	if *version {
		fmt.Fprintf(
			os.Stdout,
			"%s version %s %s/%s\n",
			os.Args[0], appVersion, runtime.GOOS, runtime.GOARCH,
		)
		os.Exit(0)
	}

	if !*usePlainProxies && !*useXORProxies && !*useXORV2Proxies {
		log.Fatalf("neither -use-plain-proxy, -use-xor-proxy or -use-xor-v2-proxies-was set")
	}

	if *sourceDomains == "" {
		log.Fatalf("no -source-domains specified")
	}

	var domains []string
	for _, d := range strings.Split(*sourceDomains, ",") {
		domains = append(domains, strings.TrimSpace(d))
	}
	domains = uniqueShuffle(domains)

	listener, err := net.Listen("tcp4", *listenAddr)
	if err != nil {
		log.Fatalf("Unable to bind port, %v\n", err)
	}
	log.Printf("Listening on %s\n", *listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection, %v\n", err)
			continue
		}

		proxies := proxy.Query(domains, *verbose)
		if len(proxies) == 0 {
			log.Printf("No proxies returned from the source domains\n")
			continue
		}

		var plainProxies []proxy.Proxy
		var xorProxies []proxy.Proxy
		var xorV2Proxies []proxy.Proxy
		for _, p := range proxies {
			if *usePlainProxies && p.Type() == typ.Plain {
				plainProxies = append(plainProxies, p)
			} else if *useXORProxies && p.Type() == typ.XOR {
				xorProxies = append(xorProxies, p)
			} else if *useXORV2Proxies && p.Type() == typ.XORV2 {
				xorV2Proxies = append(xorV2Proxies, p)
			}
		}

		var allProxies []proxy.Proxy
		if *useXORV2Proxies {
			allProxies = append(allProxies, xorV2Proxies...)
		}
		if *useXORProxies {
			allProxies = append(allProxies, xorProxies...)
		}
		if *usePlainProxies {
			allProxies = append(allProxies, plainProxies...)
		}
		if len(allProxies) == 0 {
			log.Printf("No proxies returned from the source domains match your preferred selection\n")
			continue
		}

		go handleSOCKS5Conn(conn, allProxies, *verbose)
	}
}

func uniqueShuffle(input []string) []string {
	tmp := make(map[string]bool)
	for _, i := range input {
		tmp[i] = true
	}

	out := make([]string, len(tmp))
	i := 0
	for t := range tmp {
		out[i] = t
		i++
	}

	return out
}
