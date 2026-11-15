# JWT Long Lifetime

This challenge demonstrates an API that issues access tokens valid for a whole year. The token does expire, but a leaked token stays usable for months, and there is no way to revoke a self-contained JWT before then.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: tokens are valid for one year
go run main.go serve --vulnerable=true

# fixed: tokens are valid for 15 minutes and tokens living longer than one hour are rejected
go run main.go serve --vulnerable=false
```

Generate the matching token with `go run main.go jwt --vulnerable=<true|false>`.

## How to exploit it

```bash
TOKEN=$(go run main.go jwt)
jwtop decode "$TOKEN"
# exp - iat = 365 days

jwtop crack "$TOKEN" --url http://localhost:8080
vulnapi scan curl http://localhost:8080 -H "Authorization: Bearer $TOKEN"
# vulnerable: jwt.long_lifetime (default threshold 24h, see --max-token-lifetime / --jwt-max-lifetime)
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
