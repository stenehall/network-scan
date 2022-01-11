# Network-scan

## Dependencies

Using a few external dependencies. You need nmap installed on the host and need a pushover account.

We're using external go packages for nmap, sqlite and pushover.
Should be using https://pkg.go.dev/github.com/sdomino/scribble but went with sqlite just to learn it.

## Run

```bash
make build
docker run -t network-scan -subnet 192.168.0.0/24 -pushoverRecipient "some-token" -pushoverToken "some-other-token"
```