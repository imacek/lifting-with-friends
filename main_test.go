package main

import (
	"log"
	"testing"
)

func TestParseAppleFormat(t *testing.T) {
	records, _, _ := readCsv("examples/apple-strong.csv")
	clean, err := parseAppleStrongCsvRecords(records)

	log.Println(err)

	log.Println(len(records), len(clean), clean[len(clean)-1])

	aggData := calculateExerciseTimeSeries(clean)

	log.Println(len(aggData), aggData["Deadlift"])
}

func TestParseAndroidFormat(t *testing.T) {
	records, _, _ := readCsv("examples/android-strong.csv")
	clean, err := parseAndroidStrongCsvRecords(records)

	log.Println(err)

	log.Println(len(records), len(clean), clean[len(clean)-1])

	aggData := calculateExerciseTimeSeries(clean)

	log.Println(len(aggData), aggData["Deadlift"])
}
