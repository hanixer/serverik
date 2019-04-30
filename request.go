package serverik

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type HttpRequest struct {
	Method  string
	Path    string
	Headers Headers
	Buffer  []byte
}

func (request *HttpRequest) MarshalBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("%s %s HTTP/1.1", request.Method, request.Path))
	buffer.WriteString("\r\n")
	WriteHeaders(buffer, request.Headers)
	buffer.Write(request.Buffer)
	buffer.WriteString("\r\n")

	return buffer.Bytes(), nil
}

func ParseRequestByte(b []byte) (*HttpRequest, error) {
	panic("NOT IMPLEMENTED")
}

func ReadRequest(r io.Reader) (*HttpRequest, error) {
	request := new(HttpRequest)
	request.Headers = make(Headers)
	reader := bufio.NewReader(r)

	// Request-Line   = Method SP Request-URI SP HTTP-Version CRLF
	requestLine, err := ReadWhileEmptyLines(reader)
	if err != nil {
		return nil, err
	}

	index1 := strings.Index(requestLine, " ")
	request.Method = requestLine[:index1]

	requestLine = requestLine[index1+1:]
	index2 := strings.Index(requestLine, " ")
	request.Path = requestLine[:index2]

	headers, err := ReadHeaders(reader)
	if err != nil {
		return nil, err
	}

	request.Headers = headers

	length, err := GetContentLength(request.Headers)
	if err != nil {
		return request, err
	}

	buffer := make([]byte, length)
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		log.Printf("Read only %d bytes, expected %d", n, length)
		return request, err
	}
	request.Buffer = buffer

	return request, nil
}

func SendRequestGet(url string) (HttpResponse, error) {
	var response HttpResponse

	presult := UrlParse(url)

	request := HttpRequest{}
	request.Method = "GET"
	request.Headers = Headers{
		"Host":            presult.Netloc,
		"connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate",
		"User-Agent":      "golang app for http requests",
	}
	if len(presult.Path) < 1 {
		request.Path = "/"
	} else {
		request.Path = presult.Path
	}

	conn, err := net.DialTimeout("tcp", presult.Netloc+":80", 2*time.Second)
	// conn, err := net.Dial("tcp", presult.netloc+":"+presult.scheme)
	if err != nil {
		return response, err
	}

	bytes, _ := request.MarshalBinary()
	conn.Write(bytes)

	resp, err := ReadResponse(conn)
	if err != nil {
		return response, err
	}

	return resp, nil
}
