package ipcrawler

func RunIpCrawler(apiKeys map[string]string) {
	// Read API keys file here and initalize struct for that.

	for _, agent := range agents {
		// initialize API key for agent
		key, ok := apiKeys[agent.Agent()]
		if ok {
			agent.SetApiKey(key)
		}

		agent.Crawl()
	}
}
