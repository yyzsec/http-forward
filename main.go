package main

import (
	"bytes"
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
	// get the response and decorate it with your gin framework
	resp := forward(c.Request)
	// set status code
	c.Status(resp.StatusCode)
	// show the headers to your browser
	for key, value := range resp.Header {
		c.Header(key, value[0])
	}
	// show the body in your browser
	content, err := ioutil.ReadAll(resp.Body)
	// handle err and print so this program won't crash
	if err != nil {
		log.Println(err)
	}
	if err = resp.Body.Close();err != nil {
		log.Println(err)
	}
	c.String(resp.StatusCode, string(content))
}

// forward do request to your target url and returns a Response
func forward(req *http.Request) (resp *http.Response) {
	// a very very very simple log...
	fmt.Println("Forwarding to : " + *forwardHost + "......")
	// check TSL
	// watch out, maybe https is not working well orz...
	HTTPVersion := "http"
	if req.TLS != nil {
		HTTPVersion = "https"
	}

	// warp the URL...
	// targetURL is the final url format and no more other processes
	targetURL := url.URL{
		Host: *forwardHost,
		Scheme: HTTPVersion,
		Path: req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}

	// warp the request body
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	// generate the request object
	request, err := http.NewRequest(req.Method, targetURL.String(), bytes.NewReader(content))
	// handle err and print so this program won't crash
	if err != nil {
		log.Println(err)
	}
	// warp the request header
	request.Header = req.Header
	// warp the request host
	request.Host = *forwardHost
	// do request
	client := &http.Client{}
	resp, err = client.Do(request)

	if err != nil {
		log.Println(err)
	}
	return
}

func main() {
	// parse the args from outside
	flag.Parse()
	flag.Args()

	// turn off debug mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// listen all routes
	r.Any("/*any", parseRequest)

	//log
	fmt.Println("开始运行http服务：" + *localHost)

	// print the err and exit
	if err := r.Run(*localHost);err != nil {
		log.Println(err)
		return
	}
}