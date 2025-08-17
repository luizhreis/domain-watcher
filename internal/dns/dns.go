package dns

import "net"

type dns struct{}

var _ DNS = (*dns)(nil)

func NewDNS() DNS {
	return &dns{}
}

func (d *dns) Resolve(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}
