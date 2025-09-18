package protocol

import (
	"bytes"
	"fmt"
)

type Encoder struct{}

var CRLF = "\r\n"

func (e *Encoder) encodeStringArray(sa []string) []byte {
	res := []byte(fmt.Sprintf("*%d%s", len(sa), CRLF))
	for _, s := range sa {
		res = append(res, e.Encode(s, false)...)
	}
	return res
}

func (e *Encoder) Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, CRLF))
		}
		return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF))
	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d%s", v, CRLF))
	case error:
		return []byte(fmt.Sprintf("-%s%s", v, CRLF))
	case []string:
		return e.encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range value.([][]string) {
			buf.Write(e.encodeStringArray(sa))
		}
		return []byte(fmt.Sprintf("*%d%s%s", len(value.([][]string)), CRLF, buf.Bytes()))
	default:
		return nil
	}
}
