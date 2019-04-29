package serverik

import (
	"reflect"
	"testing"
)

func TestUrlParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want UrlParseResult
	}{
		// TODO: Add test cases.
		{"full", args{"http://uinames.com"}, UrlParseResult{"http", "uinames.com", "", "", ""}},
		{"full", args{"http://uinames.com/api"}, UrlParseResult{"http", "uinames.com", "/api", "", ""}},
		{"full", args{"http://uinames.com/api?que"}, UrlParseResult{"http", "uinames.com", "/api", "que", ""}},
		{"full", args{"http://thing.com/path?query#frag"}, UrlParseResult{"http", "thing.com", "/path", "query", "frag"}},
		{"full", args{"http://thing.com/path?query"}, UrlParseResult{"http", "thing.com", "/path", "query", ""}},
		{"full-ps", args{"https://thing.com/path?query#frag"}, UrlParseResult{"https", "thing.com", "/path", "query", "frag"}},
		{"no-scheme-url", args{"/path?query#frag"}, UrlParseResult{"", "", "/path", "query", "frag"}},
		// "full", args{"http://thing.com/path?q=query#frag"}, UrlParseResult{"http", "thing.com", "/path", "query", "frag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UrlParse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UrlParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseQueries(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		// TODO: Add test cases.
		{"empty", args{""}, map[string][]string{}},
		{"one", args{"k=2"}, map[string][]string{"k": []string{"2"}}},
		{"two", args{"k=2&v=3"}, map[string][]string{"k": []string{"2"}, "v": []string{"3"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQueries(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnquote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"errorescape", args{"1%2"}, "", true},
		{"empty", args{""}, "", false},
		{"plus to space", args{"1+2"}, "1 2", false},
		{"hexescape", args{"1%202"}, "1 2", false},
		{"showplus", args{"1%2B2"}, "1+2", false},
		{"backslash", args{"%5C"}, "\\", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unquote(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Quote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"empty", args{"\\"}, "%5C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quote(tt.args.s); got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}
