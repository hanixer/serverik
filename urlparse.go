package serverik

import (
	"fmt"
	"strconv"
	"strings"
)

type UrlParseResult struct {
	Scheme   string
	Netloc   string
	Path     string
	Query    string
	Fragment string
}

func UrlParse(s string) UrlParseResult {
	result := UrlParseResult{}

	if indexSemi := strings.Index(s, "://"); indexSemi >= 0 {
		result.Scheme = s[:indexSemi]
		s = s[indexSemi+3:]
	}

	indexSl := strings.Index(s, "/")
	if indexSl >= 0 {
		result.Netloc = s[:indexSl]
		s = s[indexSl:]
	} else {
		result.Netloc = s
		return result
	}

	indexQ := strings.Index(s, "?")
	if indexQ >= 0 {
		result.Path = s[:indexQ]
		s = s[indexQ+1:]
	} else {
		result.Path = s
		return result
	}

	indexH := strings.Index(s, "#")
	if indexH >= 0 {
		result.Query = s[:indexH]
		s = s[indexH+1:]
	} else {
		result.Query = s
		return result
	}

	result.Fragment = s

	return result
}

func ParseQueries(s string) map[string][]string {
	result := make(map[string][]string)

	for len(s) > 0 {
		equal := strings.Index(s, "=")
		if equal != -1 {
			k := s[:equal]
			v := ""

			end := strings.Index(s, "&")
			if end != -1 {
				v = s[equal+1 : end]
				s = s[end+1:]
			} else {
				v = s[equal+1:]
				s = ""
			}

			if _, ok := result[k]; !ok {
				result[k] = []string{}
			}
			result[k] = append(result[k], v)
		}
	}

	return result
}

// Unquote replaces escaped characters by normal chars, for example %20 is replaced with space ' '.
// Also replaces + with space ' '.
func Unquote(s string) (string, error) {
	var builder strings.Builder

	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			if i+2 < len(s) {
				n, err := strconv.ParseInt(s[i+1:i+3], 16, 0)
				if err != nil {
					return "", err
				}
				i += 2
				builder.WriteByte(byte(n))
			} else {
				return "", fmt.Errorf("HEX SEQUENCE IS NOT FINISHED")
			}
		} else if s[i] == '+' {
			builder.WriteByte(' ')
		} else {
			builder.WriteByte(s[i])
		}
	}

	return builder.String(), nil
}

func isReserved(r byte) bool {
	switch {
	case r >= 0x41 && r <= 0x5A:
		return false
	case r >= 0x61 && r <= 0x7A:
		return false
	case r >= 0x30 && r <= 0x39:
		return false
	case r == 0x2D || r == 0x2E || r == 0x5F || r == 0x7E:
		return false
	default:
		return true
	}
}

func Quote(s string) string {
	var builder strings.Builder

	for i := 0; i < len(s); i++ {
		if isReserved(s[i]) {
			builder.WriteByte('%')
			builder.WriteString(fmt.Sprintf("%X", s[i]))
		} else {
			builder.WriteByte(s[i])
		}
	}

	return builder.String()
}
