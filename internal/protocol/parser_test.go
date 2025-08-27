package protocol

import "testing"

/*
Simple String: +<text>\r\n → returns <text>.
Bulk String: $<n>\r\n<payload>\r\n, where n is bytes, payload length exactly n.
Errors when:
  - starts with neither + nor $
  - malformed CRLF in length line
  - non-numeric length
  - payload length ≠ n
  - missing trailing \r\n
Safety: function must not panic on junk/short inputs (right now it will on some)
We’ll write tests that expose those panics so you can fix them.
*/

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, but did not panic")
		}
	}()
	f()
}

func TestParser_OK(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{"simple_basic", []byte("+OK\r\n"), "OK"},
		{"simple_empty", []byte("+\r\n"), ""},
		{"bulk_ascii", []byte("$5\r\nhello\r\n"), "hello"},
		{"bulk_empty", []byte("$0\r\n\r\n"), ""},
		// Unicode: "😊" is 4 bytes in UTF-8; length must be 4
		{"bulk_unicode", []byte("$4\r\n😊\r\n"), "😊"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := SParser(tc.in)
			mustNoErr(t, err)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
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
		{"bulk_invalidLen", []byte("$7\r\nHello\r\n")},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := SParser(tc.in)
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
			_, err := SParser(tc.in)
			mustErr(t, err)
		})
	}
}

func TestSParser_NullBulk_TODO(t *testing.T) {
	t.Skip("Decide representation for $-1 (null bulk). Currently unsupported.")
	_, _ = SParser([]byte("$-1\r\n"))
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
		_, _ = SParser([]byte(s))
	})
}
