# JWT Strong EdDSA Key

This challenge demonstrates a secure JWT implementation using EdDSA.

## How to run it

```bash
go run main.go
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the EdDSA key was generated from a predictable, low-entropy seed
go run main.go serve --vulnerable=true

# fixed: the EdDSA key was generated with a cryptographically secure random source
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
