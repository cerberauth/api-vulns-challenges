# Host Header & Routing Trust

This challenge demonstrates trusting the client-supplied `Host` header: a password-reset endpoint reflects it into the generated reset link (host header injection / poisoning), and the root route uses it to select a backend, allowing virtual-host confusion.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `POST /password-reset` — returns a reset link built from the request's host
- `GET /` — routes based on the `Host` header

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the reset link uses the raw Host header, and Host: internal-admin.local reaches the internal backend
go run main.go serve --vulnerable=true

# fixed: the reset link always uses the canonical host, Host header is ignored for routing
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
