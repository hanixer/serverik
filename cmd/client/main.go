package main

import (
	"fmt"

	s "github.com/hanixer/serverik"
)

func main() {
	fmt.Println("Now...")
	r, e := s.SendRequestGet("http://www.google.com:80")
	fmt.Println(r, e)
}
