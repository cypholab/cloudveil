package types

type Options struct {
	Network      string `help:"Network address (192.168.0.0)"`
	Netmask      int    `help:"Netmask for subnet"`
	Hostname     string `arg:"required" help:"Hostname to expose real ip address"`
	ExpectedBody string `arg:"required" help:"Expected content in response body"`
	CheckScheme  string `help:"URL scheme (http or https) to check"`
	Timeout      int    `help:"Timeout in seconds"`
	RateLimit    int    `help:"Max number of requests to send in seconds"`
}
