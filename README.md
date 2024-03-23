CloudVeil is an educational tool designed to uncover publicly exposed origin servers behind Cloudflare-protected services.

Usage:
```bash
$ ./cloudveil --help
Usage: cloudveil [--network NETWORK] [--netmask NETMASK] --hostname HOSTNAME --expectedbody EXPECTEDBODY [--timeout TIMEOUT] [--ratelimit RATELIMIT]

Options:
  --network NETWORK      Network address (192.168.0.0)
  --netmask NETMASK      Netmask for subnet [default: 24]
  --hostname HOSTNAME    Hostname to expose real ip address
  --expectedbody EXPECTEDBODY
                         Expected content in response body
  --timeout TIMEOUT      Timeout in seconds [default: 5]
  --ratelimit RATELIMIT
                         Max number of requests to send in seconds [default: 100]
  --help, -h             display this help and exit

$ ./cloudveil --network 127.0.0.0 --hostname localhost
```

Search Engine based configuration:
```yaml
api_keys:
  - name: censys
    auth:
      username: test
      password: test
```
