package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	IpAuth   = "ip auth"
	UserPass = "user pass"
	SiteUrl  = "https://api.ipify.org?format=json"
)

type Entries struct {
	ConfirmedEntries int
	FailedEntries    int
}

func DetermineProxy(proxy string) string {
	if strings.Count(proxy, ":") == 1 {
		return IpAuth
	}
	return UserPass
}

func GenerateProxyParts(proxy string) (host, port, user, pass string) {
	firstColon := strings.Index(proxy, ":")
	secondColon := FindMiddleColon(proxy)
	thirdColon := strings.LastIndex(proxy, ":")

	host = proxy[:firstColon]
	port = proxy[firstColon+1 : secondColon]
	user = proxy[secondColon+1 : thirdColon]
	pass = proxy[thirdColon+1:]
	return
}

func FindMiddleColon(proxy string) int {
	lastColonIndex := strings.LastIndex(proxy, ":")
	middleColonIndex := strings.LastIndex(proxy[:lastColonIndex], ":")
	return middleColonIndex
}

func ProxyConfig(proxy string) *url.URL {
	configProxy := "http://" + proxy
	if DetermineProxy(proxy) == UserPass {
		host, port, user, pass := GenerateProxyParts(proxy)
		configProxy = "http://" + user + ":" + pass + "@" + host + ":" + port
	}
	proxyURL, err := url.Parse(configProxy)
	if err != nil {
		log.Fatal(err)
	}
	return proxyURL
}

func CreateBasicAuth(proxy string) string {
	auth := proxy[FindMiddleColon(proxy)+1:]
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
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
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		return transport
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return transport
}

func ClientSetup(jar *cookiejar.Jar, transport *http.Transport) *http.Client {
	client := &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return client
}

func RequestSetup(siteURL *url.URL, proxy string) *http.Request {
	req, _ := http.NewRequest("GET", siteURL.String(), nil)
	//Add headers here

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

func CheckAndReport(infoStruct *Info, resp *http.Response, entriesStruct *Entries) {
	log.Println("Entered:", infoStruct.entered, "\nemail: ", infoStruct.email, "\nproxy: ", infoStruct.proxy, "\nstatus: ", resp.StatusCode)
	log.Println("*************************")
	if infoStruct.entered {
		entriesStruct.ConfirmedEntries++
		return
	}
	entriesStruct.FailedEntries++
	WriteFailure(infoStruct.email)
}

func TaskEngine(infoSlice []Info) *Entries {
	entriesStruct := Entries{}

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
		CheckAndReport(&infoStruct, resp, &entriesStruct)
	}

	return &entriesStruct
}
