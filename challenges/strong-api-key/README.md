# Strong API Key

This challenge demonstrates a secure API using a strong API key.

## How to run it

```bash
go run main.go
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: a short, predictable API key and a non-constant-time comparison are used
go run main.go serve --vulnerable=true

# fixed: a strong, high-entropy API key and a constant-time comparison are used
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
