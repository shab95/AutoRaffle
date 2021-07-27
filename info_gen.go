package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/Pallinder/go-randomdata"
)

type Info struct {
	email     string
	firstName string
	lastName  string
	zipCode   int
	phone     int
	size      float32
	proxy     string
}

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
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

	file, err := os.Open("emails")
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
	zip, _ := strconv.Atoi(randomdata.PostalCode("US"))
	return zip
}

func phoneGen() int {
	num, _ := strconv.Atoi(randomdata.PhoneNumber())
	return num
}

func sizeGen() float32 {
	sizes := []float32{7.5, 8, 8.5, 8.5, 9, 9, 9, 9, 9, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10.5, 10.5, 10.5, 10.5, 10.5, 10.5, 10.5, 11, 11, 11, 11, 11.5, 12, 12, 12.5, 13}
	randInd := rand.Intn(len(sizes) - 1)
	return sizes[randInd]
}

func InfoGen() []Info {
	emails := RetrieveEmails()
	proxies := retrieveProxies()
	var infoSlice []Info

	for i, email := range emails {
		var newInfo Info
		newInfo.email = email
		newInfo.firstName = firstNameGen()
		newInfo.lastName = lastNameGen()
		newInfo.phone = phoneGen()
		newInfo.zipCode = zipCodeGen()
		newInfo.phone = phoneGen()
		newInfo.size = sizeGen()
		newInfo.proxy = proxies[i]
		infoSlice = append(infoSlice, newInfo)
	}

	return infoSlice
}
