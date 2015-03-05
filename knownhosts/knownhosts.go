package knownhosts

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/phyber/rmsshkey/knownhost"
)

type KnownHosts struct {
	file     *os.File
	ch       chan *knownhost.KnownHost
	chClosed bool
}

func path() string {
	home := os.Getenv("HOME")
	path := filepath.Join(home, ".ssh", "known_hosts")
	return path
}

func open() (*os.File, error) {
	knownHostsPath := path()
	file, err := os.Open(knownHostsPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (k *KnownHosts) closeChannel() {
	if !k.chClosed {
		close(k.ch)
		k.chClosed = true
	}
}

func (k *KnownHosts) closeKnownHosts() {
	if k.file != nil {
		k.file.Close()
	}
}

// Returns a channel of knownhost.KnownHost that can be iterated over by
// range.
func (k *KnownHosts) Hosts() chan *knownhost.KnownHost {
	k.ch = make(chan *knownhost.KnownHost)
	scanner := bufio.NewScanner(k.file)

	go func() {
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

func (k *KnownHosts) Close() {
	k.closeChannel()
	k.closeKnownHosts()
}

func Open() (*KnownHosts, error) {
	k := &KnownHosts{}

	file, err := open()
	if err != nil {
		return nil, err
	}

	k.file = file

	return k, nil
}
