package dns

import "net"

func exists(host string, addrs []string) bool {
	for _, v := range addrs {
		if host == v {
			return true
		}
	}
	return false
}

func GetIPs(host string) ([]string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return []string{}, err
	}
	// Also append the given host to this list, for when we loop over it
	if !exists(host, addrs) {
		addrs = append(addrs, host)
	}
	return addrs, nil
}
