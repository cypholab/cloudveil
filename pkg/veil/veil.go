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
	hostname     string
	ips          []netip.Addr
	ratelimiter  *ratelimit.GlobalRateLimit
	progress     *progressbar.ProgressBar
	timeout      int
	expectedBody string
}

func NewVeil(ips []netip.Addr, body, hostname string, rlimit, timeout int) *Veil {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ratelimiter := ratelimit.NewGlobalRateLimit(rlimit)
	ratelimiter.Add(rateLimitKey)

	progress := progressbar.Default(int64(len(ips)))

	return &Veil{
		hostname:     hostname,
		ips:          ips,
		ratelimiter:  ratelimiter,
		expectedBody: body,
		timeout:      timeout,
		progress:     progress,
	}
}

func (b *Veil) Run() {
	results := make(chan Result)
	var scheme = "http"

	// TODO(implement HTTPS based requests sender)

	log.Printf("Total generated IP addresses: %v", len(b.ips))

	wg := sync.WaitGroup{}
	go func() {
		defer close(results)

		for _, ip := range b.ips {
			wg.Add(1)
			go func(ip netip.Addr) {
				defer wg.Done()

				b.ratelimiter.Take(rateLimitKey)
				b.progress.Add(1)

				url := fmt.Sprintf("%s://%s", scheme, ip)
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
