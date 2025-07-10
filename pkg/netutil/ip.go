package netutil

import (
	"net"
)

func GetLocalIPs() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip = ip.To4(); ip == nil || ip.IsLoopback() {
			continue
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

// GetCIDRHosts returns a list of IP addresses in the given CIDR
func GetCIDRHosts(cidr string) ([]net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ips := []net.IP{}
	for current := ip.Mask(ipnet.Mask); ipnet.Contains(current); inc(current) {
		ipCopy := make(net.IP, len(current))
		copy(ipCopy, current)
		ips = append(ips, ipCopy)
	}
	return ips, nil
}

// inc increases an IP address by 1
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
