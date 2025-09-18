package protocol

import (
	"fmt"
	"strconv"
)

type REPSParser struct{}

func (p *REPSParser) Parse(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	res, _, err := p.DecodeOne(data, 0)
	return res, err
}

func (p *REPSParser) DecodeOne(data []byte, pos int) (interface{}, int, error) {
	if pos >= len(data) || pos+1 >= len(data) {
		return nil, -1, fmt.Errorf("unexpected end of input")
	}
	switch data[pos] {
	case '+':
		return p.parseSimpleString(data, pos)
	case '$':
		return p.parseBulkString(data, pos)
	case ':':
		return p.parseInt(data, pos)
	case '*':
		return p.parseArray(data, pos)
	case '-':
		return p.parseErr(data, pos)
	default:
		return nil, -1, fmt.Errorf("invalid input: %q\nAt pos: %d, parts %s", data, pos, data[pos:])
	}
}

// +hello\r\n
func (p *REPSParser) parseSimpleString(data []byte, pos int) (interface{}, int, error) {
	pos++ // skips +
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid simple string format")
	}
	return string(data[start:pos]), pos + 2, nil
}

// $5\r\nhello\r\n
func (p *REPSParser) parseBulkString(data []byte, pos int) (interface{}, int, error) {
	pos++
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid bulk string format")
	}

	n, err := strconv.ParseInt(string(data[start:pos]), 10, 32)
	if err != nil {
		return nil, -1, fmt.Errorf("invalid bulk string format, invalid body size: %d", n)
	}
	start = pos + 2
	pos = start
	if n == 0 {
		return "", pos + 2, nil
	}

	if n < 0 || start+int(n) >= len(data) {
		return nil, -1, fmt.Errorf("invalid bulk string format, invalid body size: %d", n)
	}

	for i := 0; i < int(n); i++ {
		if data[pos] == '\r' {
			break
		}
		pos++
	}

	if pos >= len(data) || pos+1 >= len(data) || data[pos] != '\r' || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid bulk string format, pos %d, %s", pos, string(data))
	}
	return string(data[start:pos]), pos + 2, nil
}

// :-100\r\n
func (p *REPSParser) parseInt(data []byte, pos int) (interface{}, int, error) {
	pos++
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid INT format")
	}
	res, _ := strconv.ParseInt(string(data[start:pos]), 10, 64)
	return res, pos + 2, nil
}

// *3\r\n$5\r\nhello\r\n:10\r\n$5\r\nworld\r\n
func (p *REPSParser) parseArray(data []byte, pos int) (interface{}, int, error) {
	pos++
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid ARRAY format")
	}
	n, _ := strconv.ParseInt(string(data[start:pos]), 10, 32)
	pos += 2
	if n < 0 {
		return nil, -1, fmt.Errorf("invalid ARRAY formatnegative body size: %d", n)
	}
	res := make([]interface{}, 0)
	for n > 0 {
		element, npos, err := p.DecodeOne(data, pos)
		if err != nil {
			return nil, -1, fmt.Errorf("invalid format, %s", err)
		}
		pos = npos
		res = append(res, element)
		n--
	}
	return res, pos, nil
}

// -Key Not Found\r\n
func (p *REPSParser) parseErr(data []byte, pos int) (interface{}, int, error) {
	pos++
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos >= len(data) || pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid ERR format")
	}
	return string(data[start:pos]), pos + 2, nil
}
