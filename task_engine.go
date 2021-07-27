package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	IpAuth   = "ip auth"
	UserPass = "user pass"
	SiteUrl  = "https://www.google.com/"
)

func DetermineProxy(proxy string) string {
	if strings.Count(proxy, ":") == 1 {
		return IpAuth
	}
	return UserPass
}

func FindMiddleColon(proxy string) int {
	lastColonIndex := strings.LastIndex(proxy, ":")
	middleColonIndex := strings.LastIndex(proxy[:lastColonIndex], ":")
	return middleColonIndex
}

func ProxyConfig(proxy string) *url.URL {
	var hostPort string

	if DetermineProxy(proxy) == UserPass {
		hostPort = proxy[:FindMiddleColon(proxy)]
	} else {
		hostPort = proxy
	}

	configProxy := "http://" + hostPort
	proxyURL, err := url.Parse(configProxy)
	if err != nil {
		log.Fatal(err)
	}
	return proxyURL
}

func CreateBasicAuth(proxy string) string {
	auth := proxy[FindMiddleColon(proxy)+1:]
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	return basicAuth
}

func SiteConfig() *url.URL {
	url, err := url.Parse(SiteUrl)
	if err != nil {
		log.Fatal(err)
	}
	return url
}

func TransportSetup(proxyURL *url.URL, proxy string) *http.Transport {
	if DetermineProxy(proxy) == UserPass {
		hdr := http.Header{}
		hdr.Add("Proxy-Authorization", CreateBasicAuth(proxy))
		transport := &http.Transport{
			Proxy:              http.ProxyURL(proxyURL),
			ProxyConnectHeader: hdr,
		}
		return transport
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	// if DetermineProxy(proxy) == UserPass {
	// 	log.Println("This entry requires user-pass authentication...")
	// 	basicAuth := CreateBasicAuth(proxy)
	// 	transport.ProxyConnectHeader = http.Header{}
	// 	transport.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
	// }

	return transport
}

func ClientSetup(jar *cookiejar.Jar, transport *http.Transport) *http.Client {
	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	return client
}

func RequestSetup(siteURL *url.URL, proxy string) *http.Request {
	req, _ := http.NewRequest("GET", siteURL.String(), nil)
	//Add headers here
	if DetermineProxy(proxy) == UserPass {
		req.Header.Add("Accept-Encoding", "identity")
		req.Header.Add("Proxy-Authorization", CreateBasicAuth(proxy))
	}

	return req
}

func DoAndGet(infoStruct *Info, client *http.Client, req *http.Request, proxy string) *http.Response {
	// dump, _ := httputil.DumpRequest(req, false)
	// log.Println(string(dump))
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(proxy, err)
		infoStruct.entered = false
		log.Println(err)
		return &http.Response{}
	} else {
		infoStruct.entered = true
	}
	defer resp.Body.Close()
	return resp
}

func TaskEngine(infoSlice []Info) {
	for _, infoStruct := range infoSlice {
		//proxy url setup
		proxyURL := ProxyConfig(infoStruct.proxy)

		//site url setup
		siteURL := SiteConfig()

		//transport setup
		transport := TransportSetup(proxyURL, infoStruct.proxy)

		//cookiejar setup
		jar, _ := cookiejar.New(nil)

		//client setup
		client := ClientSetup(jar, transport)

		//request setup
		req := RequestSetup(siteURL, infoStruct.proxy)

		//DoAndGetStatus
		resp := DoAndGet(&infoStruct, client, req, infoStruct.proxy)

		log.Println("Entered:", infoStruct.entered, "\nemail: ", infoStruct.email, "\nproxy: ", infoStruct.proxy, "\nstatus: ", resp.StatusCode)
		log.Println("*************************")
	}

}
