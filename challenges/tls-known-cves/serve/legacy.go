package serve

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"net"
	"slices"
)

const (
	recordAlert     = 0x15
	recordHandshake = 0x16
	recordHeartbeat = 0x18

	sslv2ClientHello = 0x01
	sslv2ServerHello = 0x04

	extensionHeartbeat = 0x000f

	versionTLS12 = 0x0303
)

var rsaKeyExchangeSuites = []uint16{0x002f, 0x0035, 0x003c, 0x003d, 0x009c, 0x009d, 0x000a}

var leakedMemory = []byte("Cookie: session=4f1c2a9be07d43a8; Authorization: Bearer leaked-access-token\r\n")

type clientHello struct {
	version   uint16
	suites    []uint16
	heartbeat bool
}

func (h *clientHello) rsaKeyExchangeOnly() bool {
	if h.version != versionTLS12 || len(h.suites) == 0 {
		return false
	}
	for _, suite := range h.suites {
		if !slices.Contains(rsaKeyExchangeSuites, suite) {
			return false
		}
	}
	return true
}

func (h *clientHello) pick(preferred []uint16) uint16 {
	for _, suite := range preferred {
		if slices.Contains(h.suites, suite) {
			return suite
		}
	}
	return h.suites[0]
}

func peekClientHello(reader *bufio.Reader) (*clientHello, error) {
	header, err := reader.Peek(5)
	if err != nil {
		return nil, err
	}
	if header[0] != recordHandshake {
		return nil, errors.New("not a handshake record")
	}
	record, err := reader.Peek(5 + int(binary.BigEndian.Uint16(header[3:5])))
	if err != nil {
		return nil, err
	}
	msg := record[5:]
	if len(msg) < 4+2+32+1 || msg[0] != 0x01 {
		return nil, errors.New("not a client hello")
	}
	body := msg[4:]
	hello := &clientHello{version: binary.BigEndian.Uint16(body)}
	rest, ok := skipVector(body[2+32:])
	if !ok || len(rest) < 2 {
		return nil, errors.New("short client hello")
	}
	suitesLength := int(binary.BigEndian.Uint16(rest))
	if len(rest) < 2+suitesLength+1 {
		return nil, errors.New("short client hello")
	}
	for i := 0; i+1 < suitesLength; i += 2 {
		hello.suites = append(hello.suites, binary.BigEndian.Uint16(rest[2+i:]))
	}
	rest, ok = skipVector(rest[2+suitesLength:])
	if !ok || len(rest) < 2 {
		return hello, nil
	}
	extensions := rest[2:]
	for len(extensions) >= 4 {
		extType := binary.BigEndian.Uint16(extensions)
		length := int(binary.BigEndian.Uint16(extensions[2:]))
		if len(extensions) < 4+length {
			break
		}
		if extType == extensionHeartbeat {
			hello.heartbeat = true
		}
		extensions = extensions[4+length:]
	}
	return hello, nil
}

func skipVector(b []byte) ([]byte, bool) {
	if len(b) < 1 || len(b) < 1+int(b[0]) {
		return nil, false
	}
	return b[1+int(b[0]):], true
}

func discardRecord(reader *bufio.Reader) (byte, []byte, error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(reader, header); err != nil {
		return 0, nil, err
	}
	payload := make([]byte, binary.BigEndian.Uint16(header[3:5]))
	_, err := io.ReadFull(reader, payload)
	return header[0], payload, err
}

func record(contentType byte, payload []byte) []byte {
	out := []byte{contentType, 0x03, 0x03}
	out = binary.BigEndian.AppendUint16(out, uint16(len(payload)))
	return append(out, payload...)
}

func handshake(msgType byte, body []byte) []byte {
	return append([]byte{msgType, byte(len(body) >> 16), byte(len(body) >> 8), byte(len(body))}, body...)
}

func uint24(n int) []byte {
	return []byte{byte(n >> 16), byte(n >> 8), byte(n)}
}

func serverFlight(suite uint16, der []byte, extensions []byte) []byte {
	random := make([]byte, 32)
	_, _ = rand.Read(random)
	hello := binary.BigEndian.AppendUint16(nil, versionTLS12)
	hello = append(hello, random...)
	hello = append(hello, 0)
	hello = binary.BigEndian.AppendUint16(hello, suite)
	hello = append(hello, 0)
	if len(extensions) > 0 {
		hello = binary.BigEndian.AppendUint16(hello, uint16(len(extensions)))
		hello = append(hello, extensions...)
	}

	entry := append(uint24(len(der)), der...)
	flight := handshake(0x02, hello)
	flight = append(flight, handshake(0x0b, append(uint24(len(entry)), entry...))...)
	flight = append(flight, handshake(0x0e, nil)...)
	return record(recordHandshake, flight)
}

func serveSSLv2(conn net.Conn, der []byte) {
	ciphers := []byte{0x01, 0x00, 0x80, 0x07, 0x00, 0xc0}
	connectionID := make([]byte, 16)
	_, _ = rand.Read(connectionID)
	body := []byte{sslv2ServerHello, 0x00, 0x01, 0x00, 0x02}
	body = binary.BigEndian.AppendUint16(body, uint16(len(der)))
	body = binary.BigEndian.AppendUint16(body, uint16(len(ciphers)))
	body = binary.BigEndian.AppendUint16(body, uint16(len(connectionID)))
	body = append(body, der...)
	body = append(body, ciphers...)
	body = append(body, connectionID...)
	_, _ = conn.Write(append(binary.BigEndian.AppendUint16(nil, 0x8000|uint16(len(body))), body...))
}

func serveHeartbleed(conn net.Conn, reader *bufio.Reader, der []byte, hello *clientHello) {
	if _, _, err := discardRecord(reader); err != nil {
		return
	}
	extensions := []byte{0x00, 0x0f, 0x00, 0x01, 0x01}
	if _, err := conn.Write(serverFlight(hello.pick(rsaKeyExchangeSuites), der, extensions)); err != nil {
		return
	}
	for {
		contentType, payload, err := discardRecord(reader)
		if err != nil {
			return
		}
		if contentType != recordHeartbeat || len(payload) < 3 || payload[0] != 0x01 {
			continue
		}
		length := min(int(binary.BigEndian.Uint16(payload[1:3])), 16*1024-19)
		memory := bytes.Repeat(leakedMemory, length/len(leakedMemory)+1)[:length]
		response := append([]byte{0x02, byte(length >> 8), byte(length)}, memory...)
		response = append(response, make([]byte, 16)...)
		_, _ = conn.Write(record(recordHeartbeat, response))
		return
	}
}

func serveRobot(conn net.Conn, reader *bufio.Reader, key *rsa.PrivateKey, der []byte, hello *clientHello) {
	if _, _, err := discardRecord(reader); err != nil {
		return
	}
	if _, err := conn.Write(serverFlight(hello.pick(rsaKeyExchangeSuites), der, nil)); err != nil {
		return
	}
	contentType, payload, err := discardRecord(reader)
	if err != nil || contentType != recordHandshake || len(payload) < 6 || payload[0] != 0x10 {
		return
	}
	plaintext := new(big.Int).Exp(new(big.Int).SetBytes(payload[6:]), key.D, key.N).FillBytes(make([]byte, key.Size()))
	description := byte(51)
	if plaintext[0] == 0x00 && plaintext[1] == 0x02 {
		description = 20
	}
	_, _ = conn.Write(record(recordAlert, []byte{0x02, description}))
}
