# JWT HMAC Confusion via JWKS

This challenge demonstrates HMAC/RSA key confusion on an API protected by OAuth2. Access tokens are RS256 tokens obtained from the API's token endpoint (client credentials or password grant), and the public key is published as a JWKS at `/.well-known/jwks.json`, as most identity providers do. The API verifies tokens with a key callback that hands the key to whatever algorithm the token header names, so an HS256 token signed with the PEM encoding of the published public key is accepted. Nothing secret is needed: the PEM is rebuilt from the JWKS `n` and `e` values.

It is also a target for authenticated scanning: no token has to be minted in advance, the scanner fetches it from the token endpoint and the public key from the JWKS.

| Item | Value |
| --- | --- |
| Token endpoint | `POST /oauth/token` |
| Client | `challenge-client` / `challenge-client-secret` (Basic auth or form) |
| User (password grant) | `alice` / `wonderland` |
| JWKS | `GET /.well-known/jwks.json` |
| OpenID configuration | `GET /.well-known/openid-configuration` |
| OpenAPI document | `GET /openapi.json` |
| Protected operation | `GET /me` |

A new RSA key is generated each time the server starts.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: HS256 tokens are verified with the PEM public key as HMAC secret
go run main.go serve --vulnerable=true

# fixed: only RS256 tokens are accepted
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
TOKEN=$(curl -s -u challenge-client:challenge-client-secret -d grant_type=client_credentials http://localhost:8080/oauth/token | jq -r .access_token)
curl -i http://localhost:8080/me -H "Authorization: Bearer $TOKEN"
# -> 200 {"sub": "challenge-client"}

vulnapi scan openapi http://localhost:8080/openapi.json --client-id challenge-client --client-secret challenge-client-secret
# vulnerable: jwt.hmac_confusion, with the public key fetched from /.well-known/jwks.json
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
