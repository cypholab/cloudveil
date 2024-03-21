CloudVeil is an educational tool designed to uncover publicly exposed origin servers behind Cloudflare-protected services.

Usage:
```bash
$ ./cloudveil --help
Usage: cloudveil --cidr CIDR --hostname HOSTNAME [--ssl]

Options:
  --cidr CIDR            IPv4 CIDR address
  --hostname HOSTNAME    Hostname to expose real ip address
  --ssl                  Send HTTPS request
  --help, -h             display this help and exit

$ ./cloudveil --cidr 127.0.0.0/24 --hostname localhost
```