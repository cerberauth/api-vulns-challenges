# Open Proxy / SSRF via the Proxy

This challenge demonstrates a proxy fetch endpoint that acts as an open proxy: it forwards a caller-supplied URL with no allowlist, letting an attacker reach the proxy's own loopback interface or internal/private services (simulated by `/internal/secret`) that should never be reachable from the outside.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /internal/secret` — simulated internal-only service
- `GET /fetch?url=<target>` — fetches the given URL server-side and returns its body

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: GET /fetch?url=http://127.0.0.1:8080/internal/secret leaks the internal secret
go run main.go serve --vulnerable=true

# fixed: requests to loopback/private/link-local hosts are rejected with 403
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
