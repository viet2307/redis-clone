package command

import "fmt"

type Obj struct {
	Value interface{}
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func CmdSET(args []byte) ([]byte, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("invalid iput, no key or values: %q", args[0])
	}
	key, obj := args[0], Obj{}
	idx := 0
	return nil, nil
}
