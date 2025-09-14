package protocol

import (
	"fmt"
	"strconv"
)

type Command struct {
	Cmd  string
	Args []string
}

func CmdParser(body []byte) (Command, error) {
	if len(body) == 0 {
		return Command{}, fmt.Errorf("empty input")
	}
	cmd := ""
	idx := 0
	args := make([]string, 0)
	for i, c := range body {
		if c == ' ' && cmd == "" {
			cmd = string(body)[0:i]
			idx = i + 1
		} else {
			val := string(body)[idx:i]
			args = append(args, val)
			idx = i + 1
		}
	}
	return Command{
		Cmd:  cmd,
		Args: args,
	}, nil
}

func BstringParser(body []byte, pos int) (string, int, error) {
	if pos > len(body) {
		return "", -1, fmt.Errorf("invalid input\nbody: %s, pos: %d", body, pos)
	}
	pos++
	npos := pos
	for npos < len(body) && body[npos] != '\r' {
		npos++
	}
	size, err := strconv.Atoi(string(body[pos:npos]))
	if err != nil {
		return "", -1, fmt.Errorf("could not get size of body %s, start at %d and end at %d\n%q", body, pos, npos, err)
	}
	if size < 0 {
		return "", -1, fmt.Errorf("input have negetive body size %d", size)
	}
	npos += 2
	pos = npos
	for size > 0 {
		if npos+size >= len(body) || body[npos] == '\r' || body[npos] == '\n' {
			return "", -1, fmt.Errorf("invalid size of body %s, got size %d from pos %d", body, npos, err)
		}
		npos++
		size--
	}
	return string(body[pos:npos]), npos + 2, nil
}

func SParser(body []byte, pos int) (string, int, error) {
	if pos >= len(body) {
		return "", -1, fmt.Errorf("invalid input\nbody: %s, pos: %d", body, pos)
	}
	if body[pos] == '+' {
		pos++
		npos := pos
		for npos < len(body) && body[npos] != '\r' {
			npos++
		}
		if npos >= len(body) {
			return "", -1, fmt.Errorf("input don't have proper CRLF, %s", body)
		}
		return string(body[pos:npos]), npos + 2, nil
	}
	if body[pos] != '$' {
		return "", -1, fmt.Errorf("input not of type SIMPLE_STRING or BULK_STRING\nbody: %s, pos: %d", body, pos)
	}
	return BstringParser(body, pos)
}

func IntParser(body []byte, pos int) (int64, int, error) {
	if body[pos] != ':' {
		return 0, -1, fmt.Errorf("input not of type INT, body %s,pos %d", body, pos)
	}
	pos++
	npos := pos
	for npos < len(body) && body[npos] != '\r' {
		npos++
	}
	if npos >= len(body) {
		return 0, -1, fmt.Errorf("input do not have proper CRLF, body %s, pos %d", body, pos)
	}
	res, err := strconv.ParseInt(string(body[pos:npos]), 10, 64)
	fmt.Printf("body %d to INT64, pos1 %d pos2 %d", res, pos, npos)
	if err != nil {
		return 0, -1, fmt.Errorf("error convert body %s to INT64, pos %d\n%q", body, pos, err)
	}
	pos = npos + 2
	return res, pos, nil
}

func ErrorParser(body []byte, pos int) (string, int, error) {
	if len(body) == 0 {
		return "", -1, fmt.Errorf("empty input")
	}
	if body[pos] != '-' {
		return "", -1, fmt.Errorf("input not of type ERROR, body %s, pos %d", body, pos)
	}
	pos++
	end := pos
	for body[end] != '\r' {
		end++
	}
	return string(body[pos:end]), end + 2, nil
}

func ArrParser(body []byte, pos int) ([]interface{}, int, error) {
	if len(body) == 0 || pos >= len(body) {
		return nil, -1, fmt.Errorf("empty input")
	}

	if body[pos] != '*' {
		return nil, -1, fmt.Errorf("input not of type ARRAY, %s", body)
	}

	idx := pos
	for idx < len(body) && body[idx] != '\r' {
		idx++
	}
	count, err := strconv.Atoi(string(body[pos+1 : idx]))
	if err != nil {
		return nil, -1, fmt.Errorf("could not read length of input, %s\n%v", body, err)
	}
	idx += 2
	if idx > len(body) {
		return nil, -1, fmt.Errorf("out of bound at pos: %d for body: %s", idx, body)
	}

	res := make([]interface{}, 0, count)
	i := count
	for idx < len(body) && i > 0 {
		switch body[idx] {
		case '+':
			read, npos, err := SParser(body, idx)
			if err != nil {
				return nil, -1, err
			}
			res = append(res, read)
			idx = npos
		case '$':
			read, npos, err := SParser(body, idx)
			if err != nil {
				return nil, -1, err
			}
			res = append(res, read)
			idx = npos
		case ':':
			read, npos, err := IntParser(body, idx)
			if err != nil {
				return nil, -1, err
			}
			res = append(res, read)
			idx = npos
		case '*':
			read, npos, err := ArrParser(body, idx)
			if err != nil {
				return nil, -1, err
			}
			idx = npos
			res = append(res, read)
		case '-':
			read, npos, err := ErrorParser(body, idx)
			if err != nil {
				return nil, -1, err
			}
			res = append(res, read)
			idx = npos
		default:
			return nil, -1, fmt.Errorf("invalid input: %q\nAt pos: %d, parts %s", body, idx, body[idx:])
		}
		i--
	}
	if i > 0 {
		return nil, -1, fmt.Errorf("invalid input, expected %d of elements in ARRAY but got %d", count, len(res))
	}
	return res, idx, nil
}
