# Unauthenticated Endpoint

This challenge demonstrates an operation that declares bearer authentication in its OpenAPI document but returns its data to anonymous callers. It also contains the cases a scanner must not report: an operation whose security requirement is optional by design (`security: [{bearerAuth: []}, {}]`), a login redirect, and a single-page app answering `200` for every unknown path.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: GET /account returns the account without credentials
go run main.go serve --vulnerable=true

# fixed: GET /account without a valid token redirects to /login
go run main.go serve --vulnerable=false
```

| Path | Vulnerable | Fixed |
| --- | --- | --- |
| `GET /account` without credentials | `200` account JSON | `302` to `/login` |
| `GET /status` (optional security) | `200` | `200` |
| `/login`, any other path | `200` HTML | `200` HTML |

## How to exploit it

```bash
curl -i http://localhost:8080/account
# vulnerable -> 200 with the account, fixed -> 302 Location: /login

TOKEN=$(go run main.go jwt)
vulnapi scan openapi http://localhost:8080/openapi.json <<< "$TOKEN"
# vulnerable: broken_authentication.authentication_bypass on GET /account only
# fixed: no authentication bypass finding
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
