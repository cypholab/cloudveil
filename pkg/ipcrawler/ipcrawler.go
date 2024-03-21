package ipcrawler

import (
	"net"

	"github.com/c-robinson/iplib/v2"
	"github.com/cypholab/cloudveil/pkg/config"
)

func RunIpCrawler(hostname string, subnet int, apiKeys []config.ApiKey) ([]string, error) {
	ips := []string{}

	for _, agent := range agents {
		apiKey := findApiKey(agent, apiKeys)
		if apiKey != nil {
			if err := agent.SetApiKey(apiKey.Auth); err != nil {
				return nil, err
			}
		}

		ipAddrs, err := agent.CrawlIps(hostname)
		if err != nil {
			return nil, err
		}

		for _, ip := range ipAddrs {
			// extract only IPv4 addresses
			ip := iplib.NewNet4(net.ParseIP(ip), subnet)
			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}

func findApiKey(agent Agent, apiKeys []config.ApiKey) *config.ApiKey {
	for _, apiKey := range apiKeys {
		if apiKey.Name == agent.Agent() {
			return &apiKey
		}
	}

	return nil
}
