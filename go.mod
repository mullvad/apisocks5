module github.com/mullvad/apisocks5

go 1.22.3

replace github.com/mullvad/ipv6md => ./ipv6md

replace github.com/mullvad/proxy => ./proxy

require github.com/mullvad/proxy v0.0.0-00010101000000-000000000000

require (
	github.com/likexian/doh-go v0.6.4 // indirect
	github.com/likexian/gokit v0.25.13 // indirect
	github.com/mullvad/ipv6md v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
