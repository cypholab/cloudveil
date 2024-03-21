package censys

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

type ApiKey struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Censys struct {
	apiKey *ApiKey
}

const baseUrl = "https://search.censys.io"

func NewCensys() *Censys {
	return &Censys{}
}

func (s *Censys) CrawlIps(hostname string) ([]string, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/hosts/search?q=%s&virtual_hosts=EXCLUDE", baseUrl, hostname)
	results, err := s.sendRequest(apiUrl)
	if err != nil {
		return nil, err
	}

	ips := []string{}
	for _, result := range results {
		hits := result.Hits
		for _, hit := range hits {
			ips = append(ips, hit.Ip)
		}
	}

	return ips, nil
}

func (s *Censys) sendRequest(apiUrl string) ([]ResponseResult, error) {
	responses := []ResponseResult{}

	url, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.apiKey.Username, s.apiKey.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s API response body: #%v", s.Agent(), err)
	}

	result := &Response{}
	if err := json.Unmarshal(body, result); err != nil {
		log.Print(string(body))
		return nil, fmt.Errorf("unable to unmarshal %s result: #%v", s.Agent(), err)
	}

	responses = append(responses, result.Result)

	if nextEntry := result.Result.Links["next"]; nextEntry != "" {
		query := url.Query()
		query.Set("cursor", nextEntry)
		url.RawQuery = query.Encode()

		nextEntryResults, err := s.sendRequest(url.String())
		if err != nil {
			return nil, err
		}

		responses = append(responses, nextEntryResults...)
	}

	return responses, nil
}

func (s *Censys) SetApiKey(apiKey any) error {
	if err := mapstructure.Decode(apiKey, &s.apiKey); err != nil {
		return fmt.Errorf("invalid API Key type for %s #%v", s.Agent(), err)
	}

	return nil
}

func (s *Censys) Agent() string {
	return "censys"
}
