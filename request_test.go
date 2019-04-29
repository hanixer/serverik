package serverik

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadRequest(t *testing.T) {
	s := "POST / HTTP/1.1\r\n" +
		"Host: localhost:9999\r\n" +
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
		"Accept-Language: en-US,en;q=0.5\r\n" +
		"Accept-Encoding: gzip, deflate\r\n" +
		"Content-Type: application/x-www-form-urlencoded\r\n" +
		"Content-Length: 27\r\n" +
		"DNT: 1\r\n" +
		"Connection: keep-alive\r\n" +
		"Upgrade-Insecure-Requests: 1\r\n" +
		"\r\n" +
		"magic=mystery&secret=spooky\r\n"
	r := bytes.NewBufferString(s)

	t.Run("basic", func(t *testing.T) {
		got, err := ReadRequest(r)
		if err != nil {
			t.Errorf("ReadRequest() error = %v", err)
			return
		}
		if got.Headers["content-length"] != "27" {
			t.Errorf("Wrong content length. %q", got.Headers["content-length"])
		}
		body := got.Buffer
		if len(body) != 27 {
			t.Errorf("expected body len 27, got %d", len(body))
		}
		if got.Method != "POST" {
			t.Error("POST is expected")
		}
		if got.Path != "/" {
			t.Error("Wront path")
		}
	})
	t.Run("mistery", func(t *testing.T) {
		s2 := `GET / HTTP/1.1
connection: keep-alive
host: www.luxoft.com
`
		r2 := bytes.NewBufferString(s2)
		got, err := ReadRequest(r2)
		fmt.Println(got, err)
	})
	t.Run("raw", func(t *testing.T) {
		s3 := "\x47\x45\x54\x20\x2f\x20\x48\x54\x54\x50\x2f\x31\x2e\x31\x0d\x0a\x63\x6f\x6e\x6e\x65\x63\x74\x69\x6f\x6e\x3a\x20\x6b\x65\x65\x70\x2d\x61\x6c\x69\x76\x65\x0d\x0a\x68\x6f\x73\x74\x3a\x20\x77\x77\x77\x2e\x6c\x75\x78\x6f\x66\x74\x2e\x63\x6f\x6d\x0d\x0a\x0d\x0a"
		r3 := bytes.NewBufferString(s3)
		got, err := ReadRequest(r3)
		fmt.Println(got, err)
	})
}
