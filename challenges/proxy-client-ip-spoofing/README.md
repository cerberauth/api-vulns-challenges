# Client-IP & Header Trust Boundary

This challenge demonstrates a proxy trust boundary flaw: an internal-only endpoint is gated on the caller's IP address, but the vulnerable implementation trusts the client-controlled `X-Forwarded-For` / `X-Real-IP` headers instead of the real connection's remote address, allowing the loopback ACL to be bypassed by spoofing the header.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: GET /internal/admin -H "X-Forwarded-For: 127.0.0.1" bypasses the loopback-only ACL
go run main.go serve --vulnerable=true

# fixed: the ACL is enforced on the actual TCP remote address, the header is ignored
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
