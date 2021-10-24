package main

import (
	"flag"
	"net/http"
)

func main() {
	listen := flag.String("p", ":8080", "listen port")
	oldHost := flag.String("r", "matrix-client.matrix.org", "original host")
	newHost := flag.String("l", "127.0.0.1", "new host")
	flag.Parse()
	proxy := GoReverseProxy(&RProxy{
		oldHost: *oldHost,
		newHost: *newHost,
	})
	serveErr := http.ListenAndServe(*listen, proxy)
	if serveErr != nil {
		panic(serveErr)
	}
}