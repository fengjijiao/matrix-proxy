package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type RProxy struct {
	newHost string
	oldHost string
}

func GoReverseProxy(this *RProxy) *httputil.ReverseProxy {
	remote, err := url.Parse("https://"+this.oldHost)
	if err != nil {
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(request *http.Request) {
		targetQuery := remote.RawQuery
		request.URL.Scheme = remote.Scheme
		request.URL.Host = remote.Host
		request.Host = remote.Host // todo 这个是关键
		request.URL.Path, request.URL.RawPath = joinURLPath(remote, request.URL)

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
		if _, ok := request.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
		}
		request.Header.Set("Accept-Encoding", "identity")
		//log.Println("request.URL.Path：", request.URL.Path, "request.URL.RawQuery：", request.URL.RawQuery)
	}

	// 修改响应头
	proxy.ModifyResponse = func(response *http.Response) error {
		//response.Header.Add("Access-Control-Allow-Origin", "*")
		contentType := response.Header.Get("Content-Type")
		if contentType == "application/json" {
			bs, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
			err = response.Body.Close()
			if err != nil {
				return err
			}
			bs = bytes.Replace(bs, []byte(this.oldHost), []byte(this.newHost),-1)
			body := ioutil.NopCloser(bytes.NewReader(bs))
			response.Body = body
			response.ContentLength = int64(len(bs))
			response.Header.Set("Content-Length", strconv.Itoa(len(bs)))
		}
		return nil
	}

	return proxy
}

// go sdk 源码
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

// go sdk 源码
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}