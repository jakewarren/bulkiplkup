package main

import (
	"bytes"
	"io"
	"net"
	"os"

	"github.com/apex/log"
	"github.com/ogier/pflag"
)

// openStdinOrFile reads from stdin or a file based on what input the user provides
func openStdinOrFile() io.Reader {
	var err error
	r := os.Stdin
	if len(pflag.Args()) >= 1 {
		r, err = os.Open(pflag.Arg(0))
		if err != nil {
			panic(err)
		}
	}
	return r
}

func checkError(message string, err error) {
	if err != nil {
		log.Errorf("%s: %s", message, err)
	}
}

// IPRange stores an IP range
type IPRange struct {
	from, to net.IP
}

// BogonRanges is a subset of the more static IPv4 Bogon/Reserved/Private ranges.
// In other words, these ranges are such fucking bogon's that they aren't even
// out in public.
var BogonRanges = []IPRange{
	{from: net.ParseIP("0.0.0.0"), to: net.ParseIP("0.255.255.255")},
	{from: net.ParseIP("10.0.0.0"), to: net.ParseIP("10.255.255.255")},
	{from: net.ParseIP("100.64.0.0"), to: net.ParseIP("10.127.255.255")},
	{from: net.ParseIP("127.0.0.0"), to: net.ParseIP("127.255.255.255")},
	{from: net.ParseIP("169.254.0.0"), to: net.ParseIP("169.254.255.255")},
	{from: net.ParseIP("172.16.0.0"), to: net.ParseIP("172.31.255.255")},
	{from: net.ParseIP("192.0.0.0"), to: net.ParseIP("192.0.0.255")},
	{from: net.ParseIP("192.0.2.0"), to: net.ParseIP("192.0.2.255")},
	{from: net.ParseIP("192.88.99.0"), to: net.ParseIP("192.88.99.255")},
	{from: net.ParseIP("192.168.0.0"), to: net.ParseIP("192.168.255.255")},
	{from: net.ParseIP("198.18.0.0"), to: net.ParseIP("198.19.255.255")},
	{from: net.ParseIP("198.51.100.0"), to: net.ParseIP("198.51.100.255")},
	{from: net.ParseIP("203.0.113.0"), to: net.ParseIP("203.0.113.255")},
	{from: net.ParseIP("224.0.0.0"), to: net.ParseIP("239.255.255.255")},
	{from: net.ParseIP("240.0.0.0"), to: net.ParseIP("255.255.255.255")},
}

// IsRoutable returns true if the IP is a publicly routable address
func IsRoutable(ip net.IP) bool {
	for _, rr := range BogonRanges {
		if rr.Contains(ip) {
			return false
		}
	}
	return true
}

// Contains checks if a given IP is in the IPRange
func (r *IPRange) Contains(ip net.IP) bool {
	return (bytes.Compare(ip, r.from) >= 0 && bytes.Compare(ip, r.to) <= 0)
}
