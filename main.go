package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cypholab/cloudveil/pkg/config"
	"github.com/cypholab/cloudveil/pkg/ipcrawler"
	"github.com/cypholab/cloudveil/pkg/iputils"
	"github.com/cypholab/cloudveil/pkg/types"
	"github.com/cypholab/cloudveil/pkg/veil"
)

func main() {
	opts := types.Options{CheckScheme: "https", Timeout: 5, RateLimit: 100, Netmask: 24}
	arg.MustParse(&opts)

	ips := []netip.Addr{}

	if _, err := os.Stat(config.ConfigFile); err == nil {
		config, err := config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		cidrsFromAgents, err := ipcrawler.RunIpCrawler(opts.Hostname, opts.Netmask, config.ApiKeys)
		if err != nil {
			log.Fatal(err)
		}

		for _, cidr := range cidrsFromAgents {
			ipsFromCidr, err := iputils.GetIpsFromCidr(cidr)
			if err != nil {
				log.Fatal(err)
			}

			ips = append(ips, ipsFromCidr...)
		}
	}

	// parse cidr to generate ip addresses
	if opts.Network != "" {
		cidrWithSubnet := fmt.Sprintf("%s/%d", opts.Network, opts.Netmask)
		ipsFromCidr, err := iputils.GetIpsFromCidr(cidrWithSubnet)
		if err != nil {
			log.Fatal(err)
		}

		ips = append(ips, ipsFromCidr...)
	}

	log.Printf("Total generated IP addresses: %v", len(ips))

	config := &veil.Config{
		CheckScheme:  opts.CheckScheme,
		IpAddresses:  ips,
		ExpectedBody: opts.ExpectedBody,
		Hostname:     opts.Hostname,
		RateLimit:    opts.RateLimit,
		Timeout:      opts.Timeout,
	}
	veil := veil.NewVeil(config)
	veil.Run()
}
