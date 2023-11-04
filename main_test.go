package main

import (
	"log"
	"testing"
)

func TestParseAppleFormat(t *testing.T) {
	records, _ := readCsv("examples/apple-strong.csv", false)
	clean, err := parseAppleCsvStyleRecords(records)

	log.Println(err)

	log.Println(len(records), len(clean), clean[len(clean)-1])
}
