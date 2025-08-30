package protocol_test

import (
	"testing"

	"tcp-server.com/m/internal/protocol"
)

func TestParserInt_OK(t *testing.T) {
	t.Parallel()
	test := []struct {
		name string
		in   []byte
		want int64
	}{
		{"simple int32", []byte(":-100\r\n"), -100},
		{"simple int32 with +", []byte(":+1000000000000\r\n"), 1e12},
		{"simple int64", []byte(":+112\r\n"), 112},
	}

	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, pos, err := protocol.IntParser(tc.in, 0)
			mustNoErr(t, err)
			if got != tc.want {
				t.Fatalf("got %d, pos: %d, want %d", got, pos, tc.want)
			}
		})
	}
}

func TestIntParser_Err(t *testing.T) {
	t.Parallel()
	test := []struct {
		name string
		in   []byte
	}{
		{"send ping", []byte("PING\r\n")},
		{"simple_nonCRLF", []byte(":")},
		{"weird num", []byte(":+-100")},
		{"weird num2", []byte(":--100")},
		{"weird num3", []byte(":++100")},
	}

	for _, tc := range test {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, _, err := protocol.IntParser(tc.in, 0)
			mustErr(t, err)
		})
	}
}
