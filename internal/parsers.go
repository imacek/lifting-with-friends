package internal

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// Expected Strong Apple input headers:
// Date,Workout Name,Duration,Exercise Name,Set Order,Weight,Reps,Distance,Seconds,Notes,Workout Notes,RPE
func parseAppleStrongCsvRecords(records [][]string) ([]LiftingSet, error) {
	cleanRecords := make([]LiftingSet, len(records))

	for index, record := range records {
		time, err := time.ParseInLocation("2006-01-02 15:04:05", record[0], losAngelesLocation)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Time failed at row %d", index))
			return []LiftingSet{}, err
		}

		weight, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Weight failed at row %d", index))
			weight = 0
		}

		reps, err := strconv.ParseInt(record[6], 10, 32)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Reps failed at row %d", index))
			reps = 0
		}

		ls := LiftingSet{
			timestamp:    time,
			exerciseName: record[3],
			weight:       weight,
			reps:         int(reps),
		}
		ls.oneRepMax = ls.calcOneRepMax()

		cleanRecords[index] = ls
	}

	return cleanRecords, nil
}

// Expected Strong Android input headers:
// Date;Workout Name;Exercise Name;Set Order;Weight;Weight Unit;Reps;RPE;Distance;Distance Unit;Seconds;Notes;Workout Notes;Workout Duration
func parseAndroidStrongCsvRecords(records [][]string) ([]LiftingSet, error) {
	cleanRecords := make([]LiftingSet, len(records))

	for index, record := range records {
		time, err := time.ParseInLocation("2006-01-02 15:04:05", record[0], losAngelesLocation)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Time failed at row %d", index))
			return []LiftingSet{}, err
		}

		weight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Weight failed at row %d", index))
			weight = 0
		}

		reps, err := strconv.ParseInt(record[6], 10, 32)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Reps failed at row %d", index))
			reps = 0
		}

		ls := LiftingSet{
			timestamp:    time,
			exerciseName: record[2],
			weight:       weight,
			reps:         int(reps),
		}
		ls.oneRepMax = ls.calcOneRepMax()

		cleanRecords[index] = ls
	}

	return cleanRecords, nil
}

// Expected DailyStrength Android input headers:
// "Date","Workout name","Exercise","Set","Weight","Reps","Distance","Duration","Measurement unit","Notes"
func parseDailyStrengthAndroidCsvRecords(records [][]string) ([]LiftingSet, error) {
	cleanRecords := make([]LiftingSet, len(records))

	for index, record := range records {
		time, err := time.ParseInLocation("2006-01-02 15:04:05", record[0], losAngelesLocation)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Time failed at row %d", index))
			return []LiftingSet{}, err
		}

		weight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Weight failed at row %d", index))
			weight = 0
		}

		reps, err := strconv.ParseInt(record[6], 10, 32)
		if err != nil {
			log.Println(fmt.Sprintf("Parsing Reps failed at row %d", index))
			reps = 0
		}

		ls := LiftingSet{
			timestamp:    time,
			exerciseName: record[2],
			weight:       weight,
			reps:         int(reps),
		}
		ls.oneRepMax = ls.calcOneRepMax()

		cleanRecords[index] = ls
	}

	return cleanRecords, nil
}
