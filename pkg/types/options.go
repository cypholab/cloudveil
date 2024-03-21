package types

type Options struct {
	Cidr         string `help:"IPv4 CIDR address"`
	Hostname     string `arg:"required" help:"Hostname to expose real ip address"`
	ExpectedBody string `arg:"required" help:"Expected content in response body"`
	Timeout      int    `help:"Timeout in seconds"`
	RateLimit    int    `help:"Max number of requests to send in seconds"`
}
