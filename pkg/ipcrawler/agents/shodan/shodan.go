package shodan

import "fmt"

type Shodan struct {
	apiKey string
}

func NewShodan() *Shodan {
	return &Shodan{}
}

func (s *Shodan) CrawlIps(_ string) ([]string, error) {
	return nil, nil
}

func (s *Shodan) SetApiKey(apiKey any) error {
	apikey, ok := apiKey.(string)
	if !ok {
		return fmt.Errorf("invalid API Key type for %s", s.Agent())
	}

	s.apiKey = apikey

	return nil
}

func (s *Shodan) Agent() string {
	return "shodan"
}
