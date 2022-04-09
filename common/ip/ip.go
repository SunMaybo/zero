package ip

import (
	"errors"
	"net"
)

// PrivateIPv4 PrivateIPv4
func PrivateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		ipNet, ok := a.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		ip := ipNet.IP.To4()
		if IsPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func LocalHostIP() string {
	if ip, err := PrivateIPv4(); err != nil {
		return "127.0.0.1"
	} else {
		return ip.String()
	}
}

// IsPrivateIPv4 IsPrivateIPv4
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func Lower16BitPrivateIP() (uint16, error) {
	ip, err := PrivateIPv4()
	if err != nil {
		return 0, err
	}
	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}
