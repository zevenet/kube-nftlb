package types

import (
	"fmt"
)

// Protocol defines which protocols are supported.
type Protocol string

const (
	// HTTP protocol.
	HTTP = Protocol("http")
	// HTTPS protocol.
	HTTPS = Protocol("https")
)

// IPversion defines which IP versions are supported.
type IPversion byte

const (
	// IPv4 = version 4.
	IPv4 = IPversion(4)
)

// IP defines how to format the IP.
type IP []byte

var (
	// LocalHostIPv4 defines the localhost IP (IPv4).
	LocalHostIPv4 = IP{127, 0, 0, 1}
)

// Port defines which ports are supported.
type Port int

const (
	// NFTLBport defines which port is the default nftlb port.
	NFTLBport = Port(5555)
	// HTTPport defines which port is the default HTTP port.
	HTTPport = Port(80)
	// HTTPSport defines which port is the default HTTPS port.
	HTTPSport = Port(443)
)

// URL has different fields that match any regular URL.
type URL struct {
	Protocol  Protocol
	IPversion IPversion
	IP        IP
	Port      Port
	Path      string
}

// Separators inside the URL.
const (
	protocolSeparator = "://"
	ipPortSeparator   = ':'
)

func (u URL) String() string {
	var IP string
	switch u.IPversion {
	case IPv4:
		IP = u.IP.ToIPv4()
	default:
		panic(fmt.Sprintf("IP version %d not supported yet", u.IPversion))
	}
	return fmt.Sprintf("%s%s%s%c%d%s", u.Protocol, protocolSeparator, IP, ipPortSeparator, u.Port, u.Path)
}

// ToIPv4 returns a string with the referred IP formatted as IPv4.
func (ip IP) ToIPv4() string {
	return fmt.Sprintf("%v.%v.%v.%v", ip[0], ip[1], ip[2], ip[3])
}
