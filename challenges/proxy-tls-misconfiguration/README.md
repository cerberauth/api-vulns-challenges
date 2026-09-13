# TLS / Transport

This challenge demonstrates a TLS-terminating proxy with a weak configuration: obsolete protocol versions (TLS 1.0) and non-forward-secret cipher suites are accepted, and the `Strict-Transport-Security` header is never sent.

## How to run it

```bash
go run main.go serve
```

The server listens over HTTPS with an ephemeral self-signed certificate generated at startup.

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: TLS 1.0 and weak CBC/3DES ciphers are accepted, no HSTS header
go run main.go serve --vulnerable=true

# fixed: only TLS 1.2+ with forward-secret AEAD ciphers is accepted, HSTS is sent with includeSubDomains and preload
go run main.go serve --vulnerable=false
```

```bash
curl -k --tls-max 1.0 --ciphers DES-CBC3-SHA https://localhost:8080/
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
