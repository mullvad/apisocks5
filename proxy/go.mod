module github.com/mullvad/proxy

replace github.com/mullvad/ipv6md => ../ipv6md

go 1.22.4

require (
	github.com/likexian/doh-go v0.6.4
	github.com/mullvad/ipv6md v0.0.0-00010101000000-000000000000
)

require (
	github.com/likexian/gokit v0.21.11 // indirect
	golang.org/x/net v0.0.0-20191116160921-f9c825593386 // indirect
	golang.org/x/text v0.3.2 // indirect
)
