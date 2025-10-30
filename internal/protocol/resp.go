package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

type REPSParser struct{}

// Parse returns the parsed command as a slice of strings (for arrays) or a single value
func (p *REPSParser) Parse(data []byte) ([]string, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	res, _, err := p.DecodeOne(data, 0)
	if err != nil {
		return nil, err
	}

	// Convert result to string slice for command processing
	switch v := res.(type) {
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		return []string{fmt.Sprintf("%v", v)}, nil
	}
}

// DecodeOne returns the parsed value as interface{}, position, and error
func (p *REPSParser) DecodeOne(data []byte, pos int) (interface{}, int, error) {
	if pos >= len(data) {
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
	case ',':
		return p.parseFloat(data, pos)
	default:
		return nil, -1, fmt.Errorf("invalid RESP type: %q at pos: %d", data[pos], pos)
	}
}

// +hello\r\n
func (p *REPSParser) parseSimpleString(data []byte, pos int) (string, int, error) {
	pos++ // skip +
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return "", -1, fmt.Errorf("invalid simple string format")
	}
	return string(data[start:pos]), pos + 2, nil
}

// $5\r\nhello\r\n
func (p *REPSParser) parseBulkString(data []byte, pos int) (string, int, error) {
	pos++ // skip $
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return "", -1, fmt.Errorf("invalid bulk string format")
	}

	n, err := strconv.ParseInt(string(data[start:pos]), 10, 64)
	if err != nil {
		return "", -1, fmt.Errorf("invalid bulk string length: %w", err)
	}

	pos += 2 // skip \r\n after length

	if n < 0 {
		// Null bulk string
		return "", pos, nil
	}

	if n == 0 {
		// Empty string
		if pos+1 >= len(data) || data[pos] != '\r' || data[pos+1] != '\n' {
			return "", -1, fmt.Errorf("invalid empty bulk string format")
		}
		return "", pos + 2, nil
	}

	// Read exactly n bytes
	end := pos + int(n)
	if end > len(data) {
		return "", -1, fmt.Errorf("insufficient data for bulk string")
	}

	// Verify CRLF after the n bytes
	if end+1 >= len(data) || data[end] != '\r' || data[end+1] != '\n' {
		return "", -1, fmt.Errorf("missing CRLF after bulk string")
	}

	return string(data[pos:end]), end + 2, nil
}

// :-100\r\n
func (p *REPSParser) parseInt(data []byte, pos int) (int64, int, error) {
	pos++ // skip :
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return 0, -1, fmt.Errorf("invalid integer format")
	}

	val, err := strconv.ParseInt(string(data[start:pos]), 10, 64)
	if err != nil {
		return 0, -1, fmt.Errorf("invalid integer value: %w", err)
	}
	return val, pos + 2, nil
}

// *3\r\n$5\r\nhello\r\n:10\r\n$5\r\nworld\r\n
func (p *REPSParser) parseArray(data []byte, pos int) ([]string, int, error) {
	pos++ // skip *
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid array format")
	}

	n, err := strconv.ParseInt(string(data[start:pos]), 10, 64)
	if err != nil {
		return nil, -1, fmt.Errorf("invalid array length: %w", err)
	}
	pos += 2

	if n < 0 {
		// Null array
		return nil, pos, nil
	}

	if n == 0 {
		// Empty array
		return []string{}, pos, nil
	}

	result := make([]string, 0, n)
	for i := 0; i < int(n); i++ {
		element, npos, err := p.DecodeOne(data, pos)
		if err != nil {
			return nil, -1, fmt.Errorf("error parsing array element %d: %w", i, err)
		}

		// Convert element to string
		var strVal string
		switch v := element.(type) {
		case string:
			strVal = v
		case int64:
			strVal = strconv.FormatInt(v, 10)
		case float64:
			strVal = strconv.FormatFloat(v, 'f', -1, 64)
		case []string:
			// Nested array - join with spaces
			strVal = strings.Join(v, " ")
		default:
			strVal = fmt.Sprintf("%v", v)
		}

		result = append(result, strVal)
		pos = npos
	}
	return result, pos, nil
}

// -Key Not Found\r\n
func (p *REPSParser) parseErr(data []byte, pos int) (error, int, error) {
	pos++ // skip -
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return nil, -1, fmt.Errorf("invalid error format")
	}
	return fmt.Errorf("%s", string(data[start:pos])), pos + 2, nil
}

// ,1.23\r\n
func (p *REPSParser) parseFloat(data []byte, pos int) (float64, int, error) {
	pos++ // skip ,
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) || data[pos+1] != '\n' {
		return 0, -1, fmt.Errorf("invalid float format")
	}

	val, err := strconv.ParseFloat(string(data[start:pos]), 64)
	if err != nil {
		return 0, -1, fmt.Errorf("invalid float value: %w", err)
	}
	return val, pos + 2, nil
}
