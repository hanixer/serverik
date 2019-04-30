package serverik

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Headers map[string]string

func ReadWhileEmptyLines(reader *bufio.Reader) (string, error) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		} else if len(line) > 0 {
			return line, nil
		}
	}
}

func ReadHeaders(reader *bufio.Reader) (Headers, error) {
	headers := make(Headers)
	for {
		line, _ := reader.ReadString('\n')
		indexCrlf := strings.Index(line, "\r\n")

		if indexCrlf < 0 {
			return nil, fmt.Errorf("EXPECTED CRLF HEADER")
		} else if indexCrlf == 0 {
			break
		}

		indexSemi := strings.Index(line, ": ")

		k := strings.ToLower(line[:indexSemi])
		v := line[indexSemi+2 : indexCrlf]
		headers[k] = v
	}
	return headers, nil
}

func GetContentLength(headers Headers) (int, error) {
	lenStr, ok := headers["content-length"]
	if !ok {
		return 0, nil
	}

	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func WriteHeaders(w io.Writer, headers Headers) {
	for k, v := range headers {
		io.WriteString(w, fmt.Sprintf("%s: %s\r\n", k, v))
	}
}
