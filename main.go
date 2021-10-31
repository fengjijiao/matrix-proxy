package main

import (
	"flag"
	"net/http"
)

func main() {
	listen := flag.String("p", ":8080", "listen port")
	originalHost := flag.String("r", "matrix-client.matrix.org", "original host")
	proxyHost := flag.String("l", "127.0.0.1", "proxy host")
	flag.Parse()
	proxy := GoReverseProxy(&RProxy{
		oldHost: *originalHost,
		newHost: *proxyHost,
	})
	serveErr := http.ListenAndServe(*listen, proxy)
	if serveErr != nil {
		panic(serveErr)
	}
}