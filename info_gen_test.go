package main

import (
	"log"
	"reflect"
	"testing"
)

func TestRetrieveEmails(t *testing.T) {
	got := RetrieveEmails()
	want := []string{"purduemargrettqoseh20@gmail.com", "nancibethvzfd78@gmail.com", "patrickcarljtk33@gmail.com", "lavetteninaogeut52@gmail.com"}
	log.Println(got)
	if !reflect.DeepEqual(got, want) {
		t.Error("Emails lists don't match")
	}
}
