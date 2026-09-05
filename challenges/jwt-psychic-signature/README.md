# JWT Psychic Signature

This challenge demonstrates a JWT implementation that is vulnerable to the Psychic Signature attack (CVE-2022-21449).

## How to run it

```bash
go run main.go
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: uses a faulty ECDSA verification routine reproducing CVE-2022-21449 (r=0, s=0 bypasses the check)
go run main.go serve --vulnerable=true

# fixed: uses the standard library's correct ECDSA verification
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
