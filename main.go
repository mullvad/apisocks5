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
	flag.Parse()

	if *version {
		fmt.Fprintf(
			os.Stdout,
			"%s version %s %s/%s\n",
			os.Args[0], appVersion, runtime.GOOS, runtime.GOARCH,
		)
		os.Exit(0)
	}

	if !*usePlainProxies && !*useXORProxies {
		log.Fatalf("neither -use-plain-proxy or -use-xor-proxy was set")
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
		for _, p := range proxies {
			if *usePlainProxies && p.Type() == typ.Plain {
				plainProxies = append(plainProxies, p)
			} else if *useXORProxies && p.Type() == typ.XOR {
				xorProxies = append(xorProxies, p)
			}
		}
		if *usePlainProxies && len(plainProxies) == 0 && !*useXORProxies {
			log.Printf("No plain proxies returned from the source domains\n")
			continue
		}
		if *useXORProxies && len(xorProxies) == 0 && !*usePlainProxies {
			log.Printf("No XOR proxies returned from the source domains\n")
			continue
		}

		// Prefer XOR proxies before plain proxies
		go handleSOCKS5Conn(conn, append(xorProxies, plainProxies...), *verbose)
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
