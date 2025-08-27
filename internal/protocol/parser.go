package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

const CRLF string = "\r\n"

// {"bulk_noLen", []byte("$\r\nHello\r\n")},
func bulkParser(body []byte) (string, error) {
	i := 1
	for i < len(body) && body[i] != '\r' {
		i++
	}
	if i >= len(body) {
		return "", fmt.Errorf("invalid format, input don't have proper CRLF")
	}
	if i == 1 {
		return "", fmt.Errorf("invalid format, input don't have body length")
	}
	r := bytes.Runes(body[1:i])

	lenn, err := strconv.Atoi(string(r))
	if lenn == 0 {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("invalid format, input don't have proper CRLF for length of string\n%v", err)
	}

	if i+lenn+2 >= len(body) {
		return "", fmt.Errorf("invalid format, input len is out of bound")
	}
	r = bytes.Runes(body[i+2 : i+2+lenn])
	return string(r), nil
}

func SParser(body []byte) (string, error) {
	// Simple string `+<text>\r\n`
	if len(body) <= 0 {
		return "", fmt.Errorf("empty input")
	}
	if body[0] == '+' {
		i := 1
		for i < len(body) && body[i] != '\r' {
			i++
		}
		if i >= len(body) {
			return "", fmt.Errorf("invalid format, input don't have proper CRLF")
		}
		return string(body[1:i]), nil
	}

	if body[0] != '$' {
		return "", fmt.Errorf("invalid input, input not of type SIMPLE_STRING or BULK_STRING")
	}
	// Bulk string `$<len>\r\n<text>\r\n`
	return bulkParser(body)
}

func IntParser(body []byte) (int64, error) {
	return 0, nil
}

func ArrParser(body string) (interface{}, error) {
	return 0, nil
}
