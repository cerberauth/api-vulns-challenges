# JWT Missing Expiry

This challenge demonstrates an API that issues access tokens without an `exp` claim and accepts them: once leaked, such a token stays valid forever. Most JWT libraries only enforce `exp` when it is present, so the API never notices.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: tokens are issued without exp and exp is optional
go run main.go serve --vulnerable=true

# fixed: tokens expire after 15 minutes and exp is required
go run main.go serve --vulnerable=false
```

Generate the matching token with `go run main.go jwt --vulnerable=<true|false>`.

## How to exploit it

```bash
TOKEN=$(go run main.go jwt)
jwtop decode "$TOKEN"
# no exp claim

curl -i http://localhost:8080 -H "Authorization: Bearer $TOKEN"
# -> 204 No Content, today and in ten years

jwtop crack "$TOKEN" --url http://localhost:8080
vulnapi scan curl http://localhost:8080 -H "Authorization: Bearer $TOKEN"
# vulnerable: jwt.missing_expiry
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
