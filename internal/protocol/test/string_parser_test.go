package protocol_test

import (
	"bytes"
	"testing"
)

func TestParser_OK(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []byte
		want []byte
	}{
		{"simple_basic", []byte("+OK\r\n"), []byte("OK")},
		{"simple_empty", []byte("+\r\n"), []byte("")},
		{"bulk_ascii", []byte("$5\r\nhello\r\n"), []byte("hello")},
		{"bulk_empty", []byte("$0\r\n\r\n"), []byte("")},
		// Unicode: "😊" is 4 bytes in UTF-8; length must be 4
		{"bulk_unicode", []byte("$4\r\n😊\r\n"), []byte("😊")},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, pos, err := p.DecodeOne(tc.in, 0)
			mustNoErr(t, err)
			if !bytes.Equal(got, tc.want) {
				t.Fatalf("got %q, pos: %d, want %q", got, pos, tc.want)
			}
		})
	}
}

func TestParser_Err(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []byte
	}{
		{"neither_plus_nor_dollar", []byte("PING\r\n")},
		{"simple_nonCRLF", []byte("+")},
		{"simple_invalidCRLF", []byte("+")},
		{"bulk_noLen", []byte("$\r\nHello\r\n")},
		{"bulk_noBody", []byte("$5\r\n\r\n")},
		{"bulk_invalid_len", []byte("$7\r\nHello\r\n")},
		{"bulk_null", []byte("$-1\r\n\r\n")},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := p.Parse(tc.in)
			mustErr(t, err)
		})
	}
}

func TestSParser_Panics_Today_But_Should_Be_Errors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   []byte
	}{
		{"empty_input", []byte{}},
		{"simple_no_cr", []byte("+OK")},
		{"bulk_no_cr", []byte("$5x")},
		{"bulk_huge_len", []byte("$100\r\nhi\r\n")},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := p.Parse(tc.in)
			mustErr(t, err)
		})
	}
}

func TestSParser_NullBulk_TODO(t *testing.T) {
	t.Skip("Decide representation for $-1 (null bulk). Currently unsupported.")
	_, _ = p.Parse([]byte("$-1\r\n"))
}

func FuzzSParser_NoPanic(f *testing.F) {
	// Seeds: a few valid/invalid forms
	seeds := [][]byte{
		[]byte("+OK\r\n"),
		[]byte("$0\r\n\r\n"),
		[]byte("$3\r\nabc\r\n"),
		[]byte("$-1\r\n"),
		[]byte("PING\r\n"),
		[]byte(""),
	}
	for _, s := range seeds {
		f.Add(string(s))
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on input %q: %v", s, r)
			}
		}()
		_, _ = p.Parse([]byte(s))
	})
}
