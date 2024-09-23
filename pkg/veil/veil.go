package veil

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/cypholab/cloudveil/pkg/ratelimit"
	"github.com/schollz/progressbar/v3"
)

const rateLimitKey = "cloudveil"

type Veil struct {
	checkScheme  string
	hostname     string
	ips          []netip.Addr
	ratelimiter  *ratelimit.GlobalRateLimit
	progress     *progressbar.ProgressBar
	timeout      int
	expectedBody string
}

func NewVeil(config *Config) *Veil {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ratelimiter := ratelimit.NewGlobalRateLimit(config.RateLimit)
	ratelimiter.Add(rateLimitKey)

	progress := progressbar.Default(int64(len(config.IpAddresses)))

	return &Veil{
		checkScheme:  config.CheckScheme,
		hostname:     config.Hostname,
		ips:          config.IpAddresses,
		ratelimiter:  ratelimiter,
		expectedBody: config.ExpectedBody,
		timeout:      config.Timeout,
		progress:     progress,
	}
}

func (b *Veil) Run() {
	results := make(chan Result)

	wg := sync.WaitGroup{}
	go func() {
		defer close(results)

		for _, ip := range b.ips {
			wg.Add(1)
			go func(ip netip.Addr) {
				defer wg.Done()

				b.ratelimiter.Take(rateLimitKey)
				b.progress.Add(1)

				url := fmt.Sprintf("%s://%s", b.checkScheme, ip)
				r, err := b.sendRequest(url)
				if err != nil {
					return
				}

				if r.statusCode >= 200 && r.statusCode < 400 {
					if strings.Contains(r.body, b.expectedBody) {
						results <- *r
					}
				}
			}(ip)
		}

		wg.Wait()
	}()

	for result := range results {
		log.Printf("successfully found: %v", result.url)
	}
}

func (b *Veil) sendRequest(url string) (*Result, error) {
	client := http.Client{
		Timeout: time.Duration(b.timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build http request: %v", err)
	}
	req.Host = b.hostname

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to send request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}
	defer resp.Body.Close()

	respBody := string(body)
	return &Result{url: url, statusCode: resp.StatusCode, headers: resp.Header, body: respBody}, nil
}
