# JWT Claim Removal

This challenge demonstrates an API that verifies the JWT signature but never requires the registered claims: with the default parser options of most JWT libraries, `exp`, `iat` and `nbf` are only checked when present, and `iss`, `aud` and `sub` are not checked at all. A token re-signed without any of them is accepted, so a token can be made to never expire or to be replayed against another audience.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: only the signature is verified
go run main.go serve --vulnerable=true

# fixed: iss, aud, sub, exp and iat are required and checked
go run main.go serve --vulnerable=false
```

## How to exploit it

The signing key is a known dev secret (`claim-removal-dev-key`), as would be recovered by a secret-cracking tool.

```bash
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080 -H "Authorization: Bearer $TOKEN"
# -> 204 No Content

jwtop crack "$TOKEN" --url http://localhost:8080 --secret claim-removal-dev-key
# vulnerable: jwt.claim_removal, plus jwt.audience_not_validated and jwt.issuer_not_validated
# fixed: no claim removal finding

vulnapi scan curl http://localhost:8080 -H "Authorization: Bearer $TOKEN" --jwt-secret claim-removal-dev-key
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
