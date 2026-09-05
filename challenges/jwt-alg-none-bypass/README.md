# JWT None Algorithm Bypass

This challenge demonstrates a JWT implementation that is vulnerable to the 'none' algorithm bypass.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
go run main.go serve --vulnerable=true   # vulnerable: accepts tokens signed with alg "none"
go run main.go serve --vulnerable=false  # fixed: only HMAC-signed tokens are accepted, alg "none" is rejected
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
