# Path & Access Control Bypass

This challenge demonstrates a proxy-level access control bypass: `/admin` is meant to be blocked, but the vulnerable implementation does a naive, case-sensitive prefix match on the raw path, which is defeated by path normalization tricks (case variation, encoded traversal, duplicate slashes) that the backend will still resolve to the protected path.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: GET /ADMIN, GET //admin or GET /admin%2f..%2fadmin bypass the block
go run main.go serve --vulnerable=true

# fixed: the path is decoded and normalized before the ACL check, so the bypasses no longer work
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
