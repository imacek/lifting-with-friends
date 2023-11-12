package internal

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"path"
	"time"
)

type Parser func(records [][]string) ([]LiftingSet, error)

type DataSet struct {
	csvDelimiter rune
	parser       Parser
	records      [][]string
}

var strongAppleCsvHeader = `Date,Workout Name,Duration,Exercise Name,Set Order,Weight,Reps,Distance,Seconds,Notes,Workout Notes,RPE`
var strongAndroidCsvHeader = `Date;Workout Name;Exercise Name;Set Order;Weight;Weight Unit;Reps;RPE;Distance;Distance Unit;Seconds;Notes;Workout Notes;Workout Duration`
var dailyStrengthAndroidCsvHeader = `"Date","Workout name","Exercise","Set","Weight","Reps","Distance","Duration","Measurement unit","Notes"`

type UnknownDataFormat struct{}

func (e *UnknownDataFormat) Error() string { return "Unknown data format detected" }

func readCsvHeader(file *os.File) (DataSet, error) {
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return DataSet{}, scanner.Err()
	}

	switch line := scanner.Text(); line {
	case strongAppleCsvHeader:
		return DataSet{',', parseStrongAppleCsvRecords, nil}, nil
	case strongAndroidCsvHeader:
		return DataSet{';', parseStrongAndroidCsvRecords, nil}, nil
	case dailyStrengthAndroidCsvHeader:
		return DataSet{',', parseDailyStrengthAndroidCsvRecords, nil}, nil
	}

	return DataSet{}, new(UnknownDataFormat)
}

func readCsvBody(file *os.File, delimiter rune) ([][]string, error) {
	file.Seek(0, 0)

	r := csv.NewReader(file)
	r.Comma = delimiter

	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records[1:], nil
}

type LiftingSet struct {
	timestamp    time.Time
	exerciseName string
	weight       float64
	reps         int
	oneRepMax    float64
}

func (ls LiftingSet) calcOneRepMax() float64 {
	return ls.weight * (36 / (37 - float64(ls.reps)))
}

var LosAngelesTimeLocation, e = time.LoadLocation("America/Los_Angeles")

func LoadUserLiftingSets(storagePath string) map[string][]LiftingSet {
	files, err := os.ReadDir(storagePath)
	if err != nil {
		log.Fatal(err)
		return map[string][]LiftingSet{}
	}

	userLiftingSets := make(map[string][]LiftingSet, len(files))

	for _, file := range files {
		filePath := path.Join(storagePath, file.Name())

		fileHandle, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
			continue
		}

		dataSet, err := readCsvHeader(fileHandle)
		if err != nil {
			fileHandle.Close()
			log.Println(err)
			continue
		}

		dataSet.records, err = readCsvBody(fileHandle, dataSet.csvDelimiter)
		if err != nil {
			fileHandle.Close()
			log.Println(err)
			continue
		}

		fileHandle.Close()

		userLiftingSets[file.Name()], err = dataSet.parser(dataSet.records)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return userLiftingSets
}

func CanStorageAcceptFile(storagePath string, fileName string, maxStoredFileCount int) bool {
	files, err := os.ReadDir(storagePath)
	if err != nil {
		log.Fatal(err)
		return false
	}

	for _, file := range files {
		if file.Name() == fileName {
			return true
		}
	}

	return len(files) < maxStoredFileCount
}
