package main

func main() {
	//ask for email file - detect number of entries based on number of emails

	//ask for proxy file - go thru proxies until finding one less than certain ping. have to make proxy requests.
	//generate real random phone number
	//generate real random zip code
	//generate real random first and last name
	//random size

	//finalize with post request
	//check status
	//if fail, move onto next struct and keep list of failed gmails

	infoSlice := InfoGen()
	entries := TaskEngine(infoSlice)
	GetStats(entries)
	SendResults(entries)
}
