module github.com/mullvad/proxy

replace github.com/mullvad/ipv6md => ../ipv6md

go 1.22.4

require (
	github.com/likexian/doh-go v0.6.5
	github.com/mullvad/ipv6md v0.0.0-00010101000000-000000000000
)

require (
	github.com/likexian/gokit v0.25.15 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)
