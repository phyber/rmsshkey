// Package knownhosts provides a basic interface to the OpenSSH known_hosts
// file.

package knownhosts

import (
	"bufio"
	"os"

	"github.com/phyber/rmsshkey/knownhost"
)

type KnownHosts struct {
	file     *os.File
	ch       chan *knownhost.KnownHost
	chClosed bool
	done     chan struct{}
}

// closeChannel closes KnownHosts.ch and sets KnownHosts.chClosed to true.
func (k *KnownHosts) closeChannel() {
	if !k.chClosed {
		close(k.ch)
		k.chClosed = true
	}
}

// closeKnownHosts closes an open known_hosts file.
func (k *KnownHosts) closeKnownHosts() {
	if k.file != nil {
		k.file.Close()
		k.file = nil
	}
}

// openKnownHosts opens a known_hosts file.
func openKnownHosts() (*os.File, error) {
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
	k.ch = make(chan *knownhost.KnownHost)
	k.done = make(chan struct{}, 2)
	scanner := bufio.NewScanner(k.file)

	go func() {
		select {
		case <-k.done:
			return
		default:
		}

		for scanner.Scan() {
			entry := scanner.Text()

			kh, err := knownhost.New(entry)
			if err != nil {
				// TODO: Handle error
				continue
			}

			k.ch <- kh
		}

		k.closeChannel()
	}()

	return k.ch
}

// Close closes all open channels and file handles opened by the package.
// It also indicates to the "Hosts" goroutine that it should return.
func (k *KnownHosts) Close() {
	k.done <- struct{}{}
	k.closeChannel()
	k.closeKnownHosts()
}

// Open opens the OpenSSH known_hosts file and returns *KnownHosts.
func Open() (*KnownHosts, error) {
	k := &KnownHosts{}

	file, err := openKnownHosts()
	if err != nil {
		return nil, err
	}

	k.file = file

	return k, nil
}
