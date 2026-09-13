# HTTP/2 & HTTP/3 Protocol-Specific

This challenge demonstrates an HTTP/2 (h2c) server that resolves routing purely from the `:authority` pseudo-header (exposed as `r.Host`) without cross-checking it against a legacy `Host` header also present on the request, allowing a mismatch between the two to be used for routing/ACL confusion during an HTTP/2-to-HTTP/1.1 downgrade.

## How to run it

```bash
go run main.go serve
```

The server speaks HTTP/2 cleartext (h2c) as well as HTTP/1.1 on the same port.

## Endpoint

- `GET /` — routes on `:authority` (`r.Host`)

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: a mismatched ":authority" and "Host" header is accepted and routed on ":authority" alone
go run main.go serve --vulnerable=true

# fixed: a mismatch between ":authority" and "Host" is rejected with 400, per RFC 7540 §8.1.2.3
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
