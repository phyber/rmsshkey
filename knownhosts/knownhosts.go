// Package knownhosts provides a basic interface to the OpenSSH known_hosts
// file.

package knownhosts

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/phyber/rmsshkey/knownhost"
)

type empty struct{}

type KnownHosts struct {
	done       chan empty
	finished   bool
	knownHosts []string
}

// openKnownHosts opens a known_hosts file.
func openKnownHosts() (*os.File, error) {
	// path() located in platform specific file.
	knownHostsPath := path()
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Hosts returns a channel of knownhost.KnownHost which can be iterated over
// via range.
func (k *KnownHosts) Hosts() <-chan *knownhost.KnownHost {
	out := make(chan *knownhost.KnownHost)
	k.done = make(chan empty)
	k.finished = false

	go func() {
		defer close(out)
		defer k.Close()

		for _, entry := range k.knownHosts {
			kh, err := knownhost.New(entry)
			if err != nil {
				// TODO: Handle error
				continue
			}

			select {
			case out <- kh:
			case <-k.done:
				return
			}
		}
	}()
	return out
}

// Close closes all open channels and file handles opened by the package.
// It also indicates to the "Hosts" goroutine that it should return.
func (k *KnownHosts) Close() {
	if k.finished {
		return
	}
	select {
	case k.done <- empty{}:
		k.finished = true
	default:
	}
}

// Remove removes a given host from the knownHosts array.
func (k *KnownHosts) Remove(host string) error {
	for i, entry := range k.knownHosts {
		kh, err := knownhost.New(entry)
		if err != nil {
			// TODO: Handle it
			continue
		}
		if kh.Host() == host {
			k.knownHosts = append(k.knownHosts[:i], k.knownHosts[i:]...)
			break
		}
	}
	return nil
}

// Open opens the OpenSSH known_hosts file and returns *KnownHosts.
func Open() (*KnownHosts, error) {
	k := &KnownHosts{}

	file, err := openKnownHosts()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	k.knownHosts = strings.Split(string(data), "\n")

	return k, nil
}
