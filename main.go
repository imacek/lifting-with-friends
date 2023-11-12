package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	. "lifting-with-friends/internal"
	"net/http"
	"path/filepath"
	"regexp"
)

var userData map[string][4]UserExerciseTimeSeries

func reloadData(storagePathFlag string) {
	liftingSets := LoadUserLiftingSets(storagePathFlag)
	userData = CalculateUserExerciseTimeSeries(liftingSets)
}

var validUsernameRegex = regexp.MustCompile(`[^a-zA-Z0-9-_]+`)

func sanitizeUsername(str string) string {
	sanitized := validUsernameRegex.ReplaceAllString(str, "")
	if len(sanitized) < 15 {
		return sanitized
	} else {
		return sanitized[:14]
	}
}

func main() {
	portFlag := flag.Int("port", 8080, "Port the web server will listen on. Defaults to 8080.")
	storagePathFlag := flag.String("storage", "storage", "Path to the storage folder. Defaults to 'storage' folder in the working directory.")
	maxStoredFileSizeFlag := flag.Int64("storage-maxfsize", 2<<20, "The maximum allowed file size inside the storage directory. Used in conjunction with storage-maxfcount to control storage size.")
	maxStoredFileCountFlag := flag.Int("storage-maxfcount", 20, "The maximum file count that can be stored inside storage directory. Used in conjunction with storage-maxfsize to control storage size.")
	flag.Parse()

	reloadData(*storagePathFlag)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = *maxStoredFileSizeFlag
	r.LoadHTMLGlob("web/*.html")

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, userData)
	})

	r.POST("/api/upload", func(c *gin.Context) {
		username := sanitizeUsername(c.PostForm("user"))
		file, err := c.FormFile("file")

		if err != nil || len(username) == 0 {
			c.Status(http.StatusBadRequest)
			return
		}

		if !CanStorageAcceptFile(*storagePathFlag, username, *maxStoredFileCountFlag) {
			c.Status(http.StatusInsufficientStorage)
			return
		}

		err = c.SaveUploadedFile(file, filepath.Join(*storagePathFlag, username))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		reloadData(*storagePathFlag)
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Static("/assets", "web")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(fmt.Sprintf(":%d", *portFlag))
}
