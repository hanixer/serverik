package main

import (
	"fmt"
	"log"
	"net"

	. "github.com/hanixer/serverik"
)

type BookmarkServ struct {
	urlDictionary map[string]string
}

var forms = `
<!DOCTYPE html>
<title>Bookmark Server</title>
<form method="POST" action="http://localhost:8000/">
Origin URL
<br>
<textarea name="origin">www.google.com</textarea>
<br>
Shortened URL key
<br>
<textarea name="short">ge</textarea>
<br>
<button type="submit">Post it!</button>
</form>`

var wrongInput = `
<h1>Error 400</h1>
One of the forms are empty
<br>Please fill both and try again`

var wrongOrigin = `
<h1>Error 404</h1>
Requested origin page does not exist`

var wrongRecord = `
<h1>Error 400</h1>
No such record`

func checkUrl(url string) bool {
	resp, err := SendRequestGet(url)
	return err == nil && resp.Response == 200
}

func (serv *BookmarkServ) HandleGET(request *HttpRequest, conn net.Conn) {
	response := NewHttpResponse()
	response.AddHeader("Content-type", "text/html; charset=utf-8")

	if request.Path == "/" {
		response.SetResponse(200)
		response.SetStringContent(forms)
	} else {
		v, ok := serv.urlDictionary[request.Path[1:]]
		if !ok {
			response.SetResponse(400)
			response.SetStringContent(wrongRecord)
		} else {
			response.AddHeader("location", v)
			response.SetResponse(302)
		}
	}
	fmt.Printf("==>\n%q\n", string(response.ToBytes()))
	conn.Write(response.ToBytes())
	conn.Close()
}

func parseBody(buffer []byte) (string, string, bool) {
	body, _ := Unquote(string(buffer))
	qs := ParseQueries(body)
	origin, ok := qs["origin"]
	if !ok {
		return "", "", false
	}
	short, ok := qs["short"]
	if !ok {
		return "", "", false
	}
	return origin[0], short[0], len(origin[0]) > 0 && len(short[0]) > 0
}

func (serv *BookmarkServ) HandlePOST(request *HttpRequest, conn net.Conn) {
	response := NewHttpResponse()
	response.AddHeader("Content-type", "text/html; charset=utf-8")
	defer func() {
		conn.Write(response.ToBytes())
		conn.Close()
	}()

	origin, short, ok := parseBody(request.Buffer)
	log.Println("origin", origin, "short", short)
	if !ok {
		response.SetStringContent(wrongInput)
		response.SetResponse(400)
		return
	}

	if !checkUrl(origin) {
		response.SetResponse(404)
		response.SetStringContent(wrongOrigin)
		return
	}

	serv.urlDictionary[short] = origin

	response.SetResponse(303)
	response.AddHeader("Location", "/")
}

func (serv *BookmarkServ) HandleRequest(request *HttpRequest, conn net.Conn) {
	if request.Method == "GET" {
		serv.HandleGET(request, conn)
	} else if request.Method == "POST" {
		serv.HandlePOST(request, conn)
	} else {
		log.Print("Unknown method")
		response := NewHttpResponse()
		response.SetResponse(501)
		conn.Write(response.ToBytes())
		conn.Close()
	}
}
func main() {
	Serve(8000, &BookmarkServ{make(map[string]string)})
	// fmt.Println(checkUrl("https://thebest-best.narod.ru/"))

}
