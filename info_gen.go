package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/Pallinder/go-randomdata"
)

type Info struct {
	email     string
	firstName string
	lastName  string
	zipCode   int
	phone     string
	size      float32
	proxy     string
	entered   bool
}

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func ClearUserPass(proxySlice []string) []string {
	count := 0
	for i, proxy := range proxySlice {
		if strings.Count(proxy, ":") > 1 {
			proxySlice = append(proxySlice[:i-count], proxySlice[i+1-count:]...)
			count++
		}
	}
	return proxySlice
}

func RetrieveEmails() []string {
	var emails []string

	file, err := os.Open("emails")
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		emails = append(emails, email)
	}

	return emails
}

func retrieveProxies() []string {
	var proxies []string

	file, err := os.Open("proxies")
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		proxies = append(proxies, proxy)
	}

	return proxies
}

func firstNameGen() string {
	if randomdata.Boolean() {
		return randomdata.FirstName(randomdata.Male)
	}
	return randomdata.FirstName(randomdata.Female)
}

func lastNameGen() string {
	return randomdata.LastName()
}

func zipCodeGen() int {
	return 20800 + randomdata.Number(10, 99)
}

func phoneGen() string {
	if randomdata.Number(2)%2 == 0 {
		return "240" + randomdata.StringNumberExt(1, "", 7)
	}
	return "301" + randomdata.StringNumberExt(1, "", 7)
}

func sizeGen() float32 {
	sizes := []float32{7.5, 8, 8.5, 8.5, 9, 9, 9, 9, 9, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10.5, 10.5, 10.5, 10.5, 10.5, 10.5, 10.5, 11, 11, 11, 11, 11.5, 12, 12, 12.5, 13}
	randInd := rand.Intn(len(sizes) - 1)
	return sizes[randInd]
}

//Infogen: Creates a slice of structs with info for raffle entries
func InfoGen() []Info {
	emails := RetrieveEmails()
	proxies := retrieveProxies()
	//proxies = ClearUserPass(proxies)
	var infoSlice []Info

	for i, email := range emails {
		var newInfo Info
		newInfo.email = email
		newInfo.firstName = firstNameGen()
		newInfo.lastName = lastNameGen()
		newInfo.zipCode = zipCodeGen()
		newInfo.phone = phoneGen()
		newInfo.size = sizeGen()
		newInfo.proxy = proxies[i]
		infoSlice = append(infoSlice, newInfo)
	}

	return infoSlice
}
