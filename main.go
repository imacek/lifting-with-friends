package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	os "os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var uploadDirectory = "storage/"

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
	weight       float32
	reps         int
	oneRepMax    float32
}

func (ls LiftingSet) calcOneRepMax() float32 {
	return ls.weight * (36 / (37 - float32(ls.reps)))
}

func (ls LiftingSet) calcNormalizedExerciseName() string {
	return strings.Trim(strings.Replace(ls.exerciseName, "(Barbell)", "", -1), " ")
}

func parseAppleCsvStyleRecords(records [][]string) ([]LiftingSet, error) {
	// Input: Date,Workout Name,Duration,Exercise Name,Set Order,Weight,Reps,Distance,Seconds,Notes,Workout Notes,RPE
	// Output: Date,Exercise Name,Weight,Reps

	noHeaderRecords := records[1:]
	cleanRecords := make([]LiftingSet, len(noHeaderRecords))

	for index, record := range noHeaderRecords {
		time, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			return []LiftingSet{}, err
		}

		weight, err := strconv.ParseFloat(record[5], 32)
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
			weight:       float32(weight),
			reps:         int(reps),
		}
		ls.oneRepMax = ls.calcOneRepMax()
		ls.exerciseName = ls.calcNormalizedExerciseName()

		cleanRecords[index] = ls
	}

	return cleanRecords, nil
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.MaxMultipartMemory = 10 << 20 // 10 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		username := c.PostForm("user")
		log.Println(username)

		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		saveFilePath := filepath.Join(uploadDirectory, username)
		c.SaveUploadedFile(file, saveFilePath)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	router.Run(":8080")
}
