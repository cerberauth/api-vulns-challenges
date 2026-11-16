package serve

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

type peekedConn struct {
	net.Conn
	reader *bufio.Reader
}

func (c *peekedConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

type legacyListener struct {
	net.Listener
	vulnerable bool
	key        *rsa.PrivateKey
	der        []byte
	conns      chan net.Conn
}

func (l *legacyListener) Accept() (net.Conn, error) {
	conn, ok := <-l.conns
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

func (l *legacyListener) acceptLoop() {
	for {
		conn, err := l.Listener.Accept()
		if err != nil {
			close(l.conns)
			return
		}
		go l.route(conn)
	}
}

func (l *legacyListener) route(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, 64*1024)
	if l.vulnerable {
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if l.emulate(conn, reader) {
			conn.Close()
			return
		}
		_ = conn.SetReadDeadline(time.Time{})
	}
	l.conns <- &peekedConn{Conn: conn, reader: reader}
}

func (l *legacyListener) emulate(conn net.Conn, reader *bufio.Reader) bool {
	header, err := reader.Peek(3)
	if err != nil {
		return false
	}
	if header[0]&0x80 != 0 && header[2] == sslv2ClientHello {
		serveSSLv2(conn, l.der)
		return true
	}
	hello, err := peekClientHello(reader)
	if err != nil {
		return false
	}
	switch {
	case hello.heartbeat:
		serveHeartbleed(conn, reader, l.der, hello)
		return true
	case hello.rsaKeyExchangeOnly():
		serveRobot(conn, reader, l.key, l.der, hello)
		return true
	}
	return false
}

func certificate() (*rsa.PrivateKey, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	return key, der, err
}

func RunServer(port string, vulnerable bool) {
	key, der, err := certificate()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	legacy := &legacyListener{Listener: listener, vulnerable: vulnerable, key: key, der: der, conns: make(chan net.Conn)}
	go legacy.acceptLoop()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Write([]byte(`{"message": "ok"}`))
	})

	server := &http.Server{
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}},
			MinVersion:   tls.VersionTLS12,
		},
	}

	log.Println("Server started at port", port)
	log.Fatal(server.ServeTLS(legacy, "", ""))
}
