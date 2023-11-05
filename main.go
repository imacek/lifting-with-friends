package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func readCsv(fileName string, skipFirstLine bool) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	if skipFirstLine {
		if _, err := r.Read(); err != nil {
			return [][]string{}, err
		}
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
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

func (ls LiftingSet) calcNormalizedExerciseName() string {
	return strings.Trim(strings.Replace(ls.exerciseName, "(Barbell)", "", -1), " ")
}

func parseStrongCsvRecords(records [][]string) ([]LiftingSet, error) {
	// Expects input headers:
	// Date,Workout Name,Duration,Exercise Name,Set Order,Weight,Reps,Distance,Seconds,Notes,Workout Notes,RPE

	noHeaderRecords := records[1:]
	cleanRecords := make([]LiftingSet, len(noHeaderRecords))

	for index, record := range noHeaderRecords {
		time, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			return []LiftingSet{}, err
		}

		weight, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return []LiftingSet{}, err
		}

		reps, err := strconv.ParseInt(record[6], 10, 32)
		if err != nil {
			return []LiftingSet{}, err
		}

		ls := LiftingSet{
			timestamp:    time,
			exerciseName: record[3],
			weight:       weight,
			reps:         int(reps),
		}
		ls.oneRepMax = ls.calcOneRepMax()
		ls.exerciseName = ls.calcNormalizedExerciseName()

		cleanRecords[index] = ls
	}

	return cleanRecords, nil
}

type ExerciseAggData struct {
	MaxWeight    float64 `json:"maxWeight"`
	MaxOneRepMax float64 `json:"maxOneRepMax"`
	TotalVolume  float64 `json:"totalVolume"`
}

type UserExerciseTimeSeries = map[string]map[time.Time]ExerciseAggData

func calculateExerciseTimeSeries(liftingSets []LiftingSet) UserExerciseTimeSeries {
	m := make(map[string]map[time.Time]ExerciseAggData)

	for _, ls := range liftingSets {
		if _, contains := m[ls.exerciseName]; !contains {
			m[ls.exerciseName] = make(map[time.Time]ExerciseAggData)
		}
		if _, contains := m[ls.exerciseName][ls.timestamp]; !contains {
			m[ls.exerciseName][ls.timestamp] = ExerciseAggData{}
		}

		data := m[ls.exerciseName][ls.timestamp]
		m[ls.exerciseName][ls.timestamp] = ExerciseAggData{
			MaxWeight:    math.Max(data.MaxWeight, ls.weight),
			MaxOneRepMax: math.Max(data.MaxOneRepMax, ls.oneRepMax),
			TotalVolume:  data.TotalVolume + ls.weight*float64(ls.reps),
		}
	}

	return m
}

var storagePath = "storage/"
var userData map[string]UserExerciseTimeSeries

func loadStorage() {
	files, err := os.ReadDir(storagePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	userData = make(map[string]UserExerciseTimeSeries)

	for _, file := range files {
		records, err := readCsv(path.Join(storagePath, file.Name()), false)
		if err != nil {
			log.Println(err)
			continue
		}

		listingSets, err := parseStrongCsvRecords(records)
		if err != nil {
			log.Println(err)
			continue
		}

		userData[file.Name()] = calculateExerciseTimeSeries(listingSets)
	}
}

func main() {
	loadStorage()

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/api/data", func(c *gin.Context) {
		c.JSON(200, userData)
	})

	router.MaxMultipartMemory = 10 << 20 // 10 MiB
	router.POST("/upload", func(c *gin.Context) {
		username := c.PostForm("user")
		log.Println(username)

		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		saveFilePath := filepath.Join(storagePath, username)
		c.SaveUploadedFile(file, saveFilePath)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.Run()
}
