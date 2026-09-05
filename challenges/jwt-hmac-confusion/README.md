# JWT HMAC/RSA Confusion

This challenge demonstrates a JWT implementation that is vulnerable to HMAC/RSA key confusion.

## How to run it

```bash
go run main.go
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: HMAC-signed tokens are verified using the RSA public key bytes as the HMAC secret
go run main.go serve --vulnerable=true

# fixed: only RSA-signed tokens are accepted, HMAC/RSA confusion is not possible
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
