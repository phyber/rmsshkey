package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/phyber/rmsshkey/dns"
	"github.com/phyber/rmsshkey/knownhosts"
)

const (
	usage = "[OPTIONS] <HOST>"
)

var opts struct {
	Verbose     bool `short:"v" long:"verbose" description:"Show verbose output."`
	Interactive bool `short:"i" long:"interactive" description:"Confirm each key deletion."`
	DryRun      bool `short:"n" long:"dry-run" description:"Show actions but do not perform them."`
}

func printf(format string, args ...interface{}) {
	if opts.Verbose {
		fmt.Printf(format, args...)
	}
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = usage
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

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
	printf("Addresses for %q: %s\n", host, strings.Join(addrs, ", "))

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
				printf("Found host %q in known_hosts\n", addr)
			}
		}
	}
}
