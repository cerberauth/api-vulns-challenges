# JWT KID SQL Injection

This challenge demonstrates a JWT implementation that is vulnerable to SQL injection via the 'kid' header.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the kid header is concatenated directly into a SQL query
go run main.go serve --vulnerable=true

# fixed: the kid is passed as a bound query parameter
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a sample token
TOKEN=$(go run main.go jwt)

# Scan with jwtop — challenge uses table "keys", not the default "tokens"
jwtop crack "$TOKEN" --url http://localhost:8080 --kid-sql-table keys
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
