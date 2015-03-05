// Package dns provides basic DNS functions.

package dns

import "net"

// exists checks if host exists in the addrs array.
func exists(host string, addrs []string) bool {
	for _, v := range addrs {
		if host == v {
			return true
		}
	}
	return false
}

// GetIPs gets all IP addresses associated with host and returns []string
// slice containing them. The slice also has the provided host appended to it.
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
