package main

import (
	"log"
	"testing"
)

func TestParseAppleFormat(t *testing.T) {
	records, _ := readCsv("examples/apple-vinko", false)
	clean, err := parseStrongCsvRecords(records)

	log.Println(err)

	log.Println(len(records), len(clean), clean[len(clean)-1])

	aggData := calculateExerciseTimeSeries(clean)

	log.Println(len(aggData), aggData["Deadlift"])
}
