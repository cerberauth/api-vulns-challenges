package serve

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

func writeResponse(conn net.Conn, status string, body string) {
	fmt.Fprintf(conn, "HTTP/1.1 %s\r\n", status)
	fmt.Fprintf(conn, "Content-Type: application/json\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprintf(conn, "Connection: keep-alive\r\n\r\n")
	conn.Write([]byte(body))
}

func readChunkedBody(r *bufio.Reader) ([]byte, error) {
	var body []byte
	for {
		sizeLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		sizeLine = strings.TrimSpace(strings.SplitN(sizeLine, ";", 2)[0])
		size, err := strconv.ParseInt(sizeLine, 16, 64)
		if err != nil {
			return nil, err
		}
		if size == 0 {
			// consume the trailing CRLF after the terminating 0-size chunk
			r.ReadString('\n')
			return body, nil
		}
		chunk := make([]byte, size)
		if _, err := readFull(r, chunk); err != nil {
			return nil, err
		}
		body = append(body, chunk...)
		r.ReadString('\n') // consume trailing CRLF after chunk data
	}
}

func readFull(r *bufio.Reader, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := r.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// handleRequest parses and responds to one HTTP request on the connection.
// It returns false if the connection should be closed.
func handleRequest(conn net.Conn, r *bufio.Reader, vulnerable bool) bool {
	requestLine, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(requestLine) == "" {
		return false
	}

	tp := textproto.NewReader(r)
	headers, err := tp.ReadMIMEHeader()
	if err != nil {
		return false
	}

	contentLength := headers.Get("Content-Length")
	transferEncoding := headers.Get("Transfer-Encoding")

	if !vulnerable && contentLength != "" && transferEncoding != "" {
		// fixed: per RFC 7230 §3.3.3, a request with both Content-Length and
		// Transfer-Encoding is ambiguous and must be rejected outright,
		// closing the connection to remove any chance of a desync
		writeResponse(conn, "400 Bad Request", `{"error": "ambiguous Content-Length/Transfer-Encoding"}`)
		return false
	}

	// vulnerable: when both headers are present, Transfer-Encoding takes
	// priority and the request boundary is determined by the chunked body,
	// exactly the "TE.CL" desync pattern - any bytes appended by the
	// attacker after the terminating chunk are left buffered on the
	// connection and get interpreted as the start of a smuggled request
	if strings.EqualFold(transferEncoding, "chunked") {
		if _, err := readChunkedBody(r); err != nil {
			return false
		}
	} else if contentLength != "" {
		n, err := strconv.Atoi(contentLength)
		if err != nil {
			return false
		}
		buf := make([]byte, n)
		if _, err := readFull(r, buf); err != nil {
			return false
		}
	}

	writeResponse(conn, "200 OK", `{"message": "request processed"}`)
	return true
}

func RunServer(port string, vulnerable bool) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			r := bufio.NewReader(c)
			for {
				if !handleRequest(c, r, vulnerable) {
					return
				}
				if r.Buffered() == 0 {
					return
				}
				// vulnerable: leftover buffered bytes (the smuggled
				// request hidden in the previous request's body) are
				// processed and answered as if they were a legitimate
				// pipelined request on the same connection
			}
		}(conn)
	}
}
