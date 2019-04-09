package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

var config struct {
	server   string
	port     int
	password string
	method   string
}

func httpSocks5(uri string) *http.Client {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	host, _, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		if parsedURL.Scheme == "https" {
			host = net.JoinHostPort(parsedURL.Host, "443")
		} else {
			host = net.JoinHostPort(parsedURL.Host, "80")
		}
	} else {
		host = parsedURL.Host
	}
	rawAddr, err := ss.RawAddr(host)
	if err != nil {
		log.Fatal(err)
	}
	serverAddr := net.JoinHostPort(config.server, strconv.Itoa(config.port))
	cipher, err := ss.NewCipher(config.method, config.password)
	if err != nil {
		log.Fatal(err)
	}
	dailFunc := func(network, addr string) (net.Conn, error) {
		return ss.DialWithRawAddr(rawAddr, serverAddr, cipher.Copy())
	}
	tr := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	tr.Dial = dailFunc
	return &http.Client{Transport: tr}
}
func main() {
	config.method = "aes-256-cfb" // method
	config.password = "服務端密碼"
	config.port = 8989 // your port
	config.server = "服務端ip"
	var uri string = "https://www.google.com.hk/?gws_rd=ssl" //轉發的請求
	client := httpSocks5(uri)
	resp, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
