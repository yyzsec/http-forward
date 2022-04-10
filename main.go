package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var localHost = flag.String("l", "0.0.0.0:9001", "需要绑定本地的http服务HOST")
var forwardHost = flag.String("r", "www.baidu.com", "需要转发的http服务HOST")

func parseRequest(c *gin.Context) {
	resp := forward(c.Request)
	c.Status(resp.StatusCode)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err = resp.Body.Close();err != nil {
		panic(err)
	}
	for key, value := range resp.Header {
		c.Header(key, value[0])
	}
	c.String(resp.StatusCode, string(content))
}

func forward(req *http.Request) (resp *http.Response) {
	fmt.Println("Forwarding to : " + *forwardHost + "......")
	HTTPVersion := "http"
	if req.TLS != nil {
		HTTPVersion = "https"
	}

	targetURL := url.URL{
		Host: *forwardHost,
		Scheme: HTTPVersion,
		Path: req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}
	request := http.Request{
		URL: &targetURL,
		Host: *forwardHost,
		Method: req.Method,
		Header: req.Header,
		Body: req.Body,
	}
	client := &http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		log.Println(err)
	}
	return
}

func main() {
	flag.Parse()
	flag.Args()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Any("/*any", parseRequest)
	fmt.Println("开始运行http服务：" + *localHost)
	if err := r.Run(*localHost);err != nil {
		return
	}
}