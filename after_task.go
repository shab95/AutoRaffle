package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetStats(entriesStruct *Entries) {
	log.Println("Successful Entries:", entriesStruct.ConfirmedEntries, "\nFailed Entries:", entriesStruct.FailedEntries)
}

//WriteFailure : Marks down accounts that have successfully changed password
func WriteFailure(email string) {
	f, err := os.OpenFile("Failed List", os.O_APPEND, 0644)
	f.WriteString(email + "\n")
	if err != nil {
		log.Fatalln("Did not write to file")
	}
}

//SendResults: Sends Results Webhook to Discord
func SendResults(entryStruct *Entries) {
	cli := http.Client{}
	url := "https://discord.com/api/webhooks/839362986148757504/DtpS2fpkbV4RgaI9Yp6iDm5oc6Pw2k73levETI0ngaWnGAZ0Hq-OGAmwFcdWa9O41ok4"
	var jsonStr = []byte(`{
		"content": null,
		"embeds": [
		  {
			"title": "Raffle Entry",
			"color": 5814783,
			"fields": [
				{
				  "name": "Successful Entries",
				  "value": "` + strconv.Itoa(entryStruct.ConfirmedEntries) + `"
				},
				{
					"name": "Failed Entries:",
					"value": "` + strconv.Itoa(entryStruct.FailedEntries) + `"
				}
			]
		  }
		]
	}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("Discord response", resp.Status)
}
