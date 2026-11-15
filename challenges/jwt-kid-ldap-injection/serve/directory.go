package serve

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Entry map[string]string

type filter func(Entry) bool

type parser struct {
	input string
	pos   int
}

func ParseFilter(input string) (func(Entry) bool, error) {
	p := &parser{input: input}
	return p.parse()
}

func (p *parser) parse() (filter, error) {
	if p.pos >= len(p.input) || p.input[p.pos] != '(' {
		return nil, errors.New("expected (")
	}
	p.pos++
	if p.pos >= len(p.input) {
		return nil, errors.New("unexpected end of filter")
	}

	var f filter
	var err error
	switch p.input[p.pos] {
	case '&', '|':
		op := p.input[p.pos]
		p.pos++
		var children []filter
		for p.pos < len(p.input) && p.input[p.pos] == '(' {
			child, err := p.parse()
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		}
		f = combine(op, children)
	case '!':
		p.pos++
		child, err := p.parse()
		if err != nil {
			return nil, err
		}
		f = func(e Entry) bool { return !child(e) }
	default:
		f, err = p.parseItem()
		if err != nil {
			return nil, err
		}
	}

	if p.pos >= len(p.input) || p.input[p.pos] != ')' {
		return nil, errors.New("expected )")
	}
	p.pos++
	return f, nil
}

func combine(op byte, children []filter) filter {
	return func(e Entry) bool {
		for _, child := range children {
			if child(e) != (op == '&') {
				return op != '&'
			}
		}
		return op == '&'
	}
}

func (p *parser) parseItem() (filter, error) {
	end := strings.IndexByte(p.input[p.pos:], ')')
	if end < 0 {
		return nil, errors.New("unterminated item")
	}
	item := p.input[p.pos : p.pos+end]
	p.pos += end

	attr, value, ok := strings.Cut(item, "=")
	if !ok || attr == "" {
		return nil, errors.New("invalid item")
	}
	attr = strings.ToLower(attr)

	if value == "*" {
		return func(e Entry) bool { _, ok := e[attr]; return ok }, nil
	}

	parts := strings.Split(value, "*")
	for i, part := range parts {
		decoded, err := unescape(part)
		if err != nil {
			return nil, err
		}
		parts[i] = strings.ToLower(decoded)
	}
	return func(e Entry) bool {
		v, ok := e[attr]
		return ok && matchParts(strings.ToLower(v), parts)
	}, nil
}

func matchParts(value string, parts []string) bool {
	if len(parts) == 1 {
		return value == parts[0]
	}
	if !strings.HasPrefix(value, parts[0]) {
		return false
	}
	value = value[len(parts[0]):]
	for _, part := range parts[1 : len(parts)-1] {
		i := strings.Index(value, part)
		if i < 0 {
			return false
		}
		value = value[i+len(part):]
	}
	return strings.HasSuffix(value, parts[len(parts)-1])
}

func unescape(s string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			b.WriteByte(s[i])
			continue
		}
		if i+2 >= len(s) {
			return "", errors.New("invalid escape")
		}
		decoded, err := hex.DecodeString(s[i+1 : i+3])
		if err != nil {
			return "", err
		}
		b.Write(decoded)
		i += 2
	}
	return b.String(), nil
}

func EscapeFilterValue(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '*', '(', ')', '\\', 0:
			b.WriteString("\\" + hex.EncodeToString([]byte{c}))
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}
