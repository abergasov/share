package utils

import (
	"fmt"
	"net"
	"strings"
)

func NewListener() (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tcp address: %w", err)
	}
	return net.ListenTCP("tcp", addr)
}

func GetInterface() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}
	// handle err
	result := make([]string, 0)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses: %w", err)
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if ip.IsLoopback() || ip.To4() == nil || filterVirtual(ip.String()) {
				continue
			}
			result = append(result, ip.String())
		}
	}
	return result, nil
}

func filterVirtual(addr string) bool {
	if strings.HasPrefix(addr, "172.") { // docker prefix
		return true
	}
	return false
}
