package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/phyber/rmsshkey/dns"
	"github.com/phyber/rmsshkey/knownhosts"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: knownhosts [host]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please specify a host")
		os.Exit(1)
	}

	host := args[0]

	addrs, err := dns.GetIPs(host)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	khs, err := knownhosts.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer khs.Close()

	for kh := range khs.Hosts() {
		for _, addr := range addrs {
			var found bool
			if found, err = kh.Equals(addr); err != nil {
				fmt.Printf("Err looking for %q: %s\n", addr, err)
			}
			if found {
				fmt.Printf("Found host %q in known_hosts\n", addr)
			}
		}
	}
}
