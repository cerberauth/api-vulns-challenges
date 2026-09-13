# Rate Limiting & DoS Resilience

This challenge demonstrates a rate limiter keyed on the client-controlled `X-Forwarded-For` header instead of the real connection address, letting an attacker bypass the throttle on a login endpoint simply by rotating the header value on every request.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `POST /api/login` — limited to 5 requests per 10 seconds per client

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: sending a different X-Forwarded-For on each request resets the limit
go run main.go serve --vulnerable=true

# fixed: the limit is enforced on the actual TCP remote address, the header is ignored
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
