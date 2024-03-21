package main

import (
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
	opts := types.Options{Timeout: 5, RateLimit: 100}
	arg.MustParse(&opts)

	ips := []netip.Addr{}

	if _, err := os.Stat(config.ConfigFile); err == nil {
		config, err := config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		apiKeys := map[string]string{}
		for _, k := range config.ApiKeys {
			apiKeys[k.Name] = k.Key
		}

		ipcrawler.RunIpCrawler(apiKeys)
	}

	// parse cidr to generate ip addresses
	if opts.Cidr != "" {
		ipsFromCidr, err := iputils.GetIpsFromCidr(opts.Cidr)
		if err != nil {
			log.Fatal(err)
		}

		ips = append(ips, ipsFromCidr...)
	}

	veil := veil.NewVeil(ips, opts.ExpectedBody, opts.Hostname, opts.RateLimit, opts.Timeout)
	veil.Run()
}
