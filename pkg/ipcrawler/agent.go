package ipcrawler

import (
	"github.com/cypholab/cloudveil/pkg/ipcrawler/agents/censys"
	"github.com/cypholab/cloudveil/pkg/ipcrawler/agents/shodan"
)

type Agent interface {
	CrawlIps(hostname string) ([]string, error)
	SetApiKey(any) error
	Agent() string
}

var agents = []Agent{shodan.NewShodan(), censys.NewCensys()}
