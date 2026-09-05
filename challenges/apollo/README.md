# Apollo GraphQL Server

This challenge demonstrates an Apollo GraphQL server.

## How to run it

```bash
npm install && npm start
```

## Modes

The server supports two modes, toggled with the `VULNERABLE` environment variable (defaults to `true`):

```bash
VULNERABLE=true npm start   # vulnerable: introspection enabled, CORS allows any origin
VULNERABLE=false npm start  # fixed: introspection disabled, CORS restricted to a trusted origin
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
