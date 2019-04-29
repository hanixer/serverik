package serverik

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type HttpResponse struct {
	Response int
	Headers  Headers
	Buffer   []byte
}

func NewHttpResponse() *HttpResponse {
	v := &HttpResponse{}
	v.Headers = make(Headers)
	return v
}

func (resp *HttpResponse) SetResponse(code int) {
	resp.Response = code
}

func (resp *HttpResponse) AddHeader(k string, v string) {
	resp.Headers[k] = v
}

func (resp *HttpResponse) ToBytes() []byte {
	builder := new(strings.Builder)
	builder.WriteString("HTTP/1.1 ")
	builder.WriteString(fmt.Sprintf("%d", resp.Response))
	builder.WriteString(" OK\r\n")
	WriteHeaders(builder, resp.Headers)
	builder.WriteString("\r\n")
	builder.Write(resp.Buffer)
	return []byte(builder.String())
}

func isChunkedEncoding(headers Headers) bool {
	v, ok := headers["transfer-encoding"]
	return ok && v == "chunked"
}

func readRestrictedBody(r *bufio.Reader, headers Headers) ([]byte, error) {
	length, err := GetContentLength(headers)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, length)
	_, errBuf := io.ReadFull(r, buffer)
	if errBuf != nil {
		return nil, errBuf
	}

	return buffer, nil
}

func readChunkedBody(r *bufio.Reader, headers Headers) ([]byte, error) {
	buffer := new(bytes.Buffer)
	for {
		b, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		hexLength := strings.TrimSuffix(string(b), "\r\n")
		if len(hexLength) == 0 {
			break
		}
		count, err := strconv.ParseInt(hexLength, 16, 0)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			// finished
			break
		}
		b = make([]byte, count)
		_, errFul := io.ReadFull(r, b)
		if errFul != nil {
			return nil, err
		}
		buffer.Write(b)
	}
	return buffer.Bytes(), nil
}

func readResponseBody(r *bufio.Reader, headers Headers) ([]byte, error) {
	if isChunkedEncoding(headers) {
		return readChunkedBody(r, headers)
	}

	return readRestrictedBody(r, headers)
}

func ReadResponse(r io.Reader) (HttpResponse, error) {
	var response HttpResponse
	reader := bufio.NewReader(r)
	response.Headers = make(Headers)
	statusLine := ReadWhileEmptyLines(reader)

	if len(statusLine) < 1 {
		return response, fmt.Errorf("empty status line")
	}

	index1 := strings.Index(statusLine, " ")
	statusLine = statusLine[index1+1:]
	index2 := strings.Index(statusLine, " ")
	respCode, _ := strconv.Atoi(statusLine[:index2])
	response.SetResponse(respCode)

	headers, err := ReadHeaders(reader)
	if err != nil {
		return response, err
	}
	response.Headers = headers

	buffer, err := readResponseBody(reader, response.Headers)
	if err != nil {
		return response, err
	}
	response.Buffer = buffer

	return response, nil
}

func (response *HttpResponse) SetStringContent(s string) {
	buffer := new(bytes.Buffer)
	buffer.WriteString(s)
	response.Buffer = buffer.Bytes()
	response.AddHeader("content-length", fmt.Sprintf("%d", len(response.Buffer)))
}
