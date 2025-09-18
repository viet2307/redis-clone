package protocol_test

import (
	"reflect"
	"testing"
)

func TestArrayParser_OK(t *testing.T) {
	t.Parallel()
	test := []struct {
		name string
		in   []byte
		want []interface{}
	}{
		{"mixed string array", []byte("*2\r\n$5\r\nHello\r\n+World!!!\r\n"), []interface{}{"Hello", "World!!!"}},
		{"int array", []byte("*3\r\n:100\r\n:-100\r\n:10000000000\r\n"), []interface{}{int64(100), int64(-100), int64(1e10)}},
		{"mix all array", []byte("*3\r\n$5\r\nHello\r\n:100\r\n$3\r\nBye\r\n"), []interface{}{"Hello", int64(100), "Bye"}},
		{"nested string array", []byte("*2\r\n$5\r\nHello\r\n*1\r\n+bye\r\n"), []interface{}{
			"Hello", []interface{}{"bye"},
		}},
		{"nested int array", []byte("*2\r\n:10\r\n*1\r\n:5\r\n"), []interface{}{
			int64(10), []interface{}{int64(5)},
		}},
		{"nested mix array", []byte("*2\r\n*1\r\n+hello\r\n*1\r\n:10\r\n"), []interface{}{
			[]interface{}{"hello"}, []interface{}{int64(10)},
		}},
	}
	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := p.Parse(tc.in)
			mustNoErr(t, err)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
