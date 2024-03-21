package ipcrawler

import "github.com/cypholab/cloudveil/pkg/ipcrawler/agents/shodan"

type Agent interface {
	Crawl()
	SetApiKey(string)
	Agent() string
}

var agents = []Agent{shodan.NewShodan()}
