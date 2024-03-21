package shodan

type Shodan struct {
	apiKey string
}

func NewShodan() *Shodan {
	return &Shodan{}
}

func (s *Shodan) Crawl() {}

func (s *Shodan) SetApiKey(apikey string) {
	s.apiKey = apikey
}

func (s *Shodan) Agent() string {
	return "shodan"
}
