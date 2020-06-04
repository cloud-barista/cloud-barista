package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// ===== [ Constants and Variables ] =====

var (
	task      config
	secretKey string
)

// ===== [ Types ] =====

// hmac - description
type config struct {
	SecretKey string `mapstructure:"secret_key"`
	AccessKey string `mapstructure:"access_key"`
	Duration  string `mapstructure:"duration"`
	Timestamp string
	Token     string
	Message   string
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// getTime - description
func getTime() time.Time {
	return time.Now().UTC()
}

// parseDuration - description
func parseDuration(limitTime string) time.Duration {
	duration, err := time.ParseDuration(limitTime)
	if err != nil {
		return 0
	}

	return duration
}

// checkDuration - description
func checkDuration(checkTime time.Time, timestamp string, duration time.Duration) bool {
	ts, err := time.Parse(time.UnixDate, timestamp)
	if err != nil {
		return false
	}

	ts = ts.Add(duration)

	log.Printf("current Time: %v", checkTime)
	log.Printf("Durable timestamp: %v", ts)
	log.Printf("Difference: %v", checkTime.Sub(ts))

	return ts.Sub(checkTime) >= 0
}

// makeToken - description
func makeToken(task config) []byte {
	data := task.Duration + "^" + task.Timestamp + "^" + task.AccessKey

	// Create HMAC
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	return h.Sum(nil)
}

// getHMACToken - description
func getHMACToken(task config) config {
	token := append(makeToken(task), []byte("||"+task.Timestamp+"|"+task.Duration+"|"+task.AccessKey)...)
	task.Token = hex.EncodeToString(token)

	return task
}

// getTokenData - description
func getTokenData(token string) [][]byte {
	log.Printf("received token: [%v]", token)

	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return [][]byte{}
	}

	sep := []byte("||")

	data := bytes.Split(tokenBytes, sep)

	return data
}

// getTask - description
func getConfigInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, task)
	}
}

// createToken - description
func createToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		var newTask config
		c.Bind(&newTask)

		newTask.Timestamp = getTime().Format(time.UnixDate)

		return c.JSON(http.StatusOK, getHMACToken(newTask))
	}
}

// validateToken - description
func validateToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		var newTask config
		c.Bind(&newTask)

		var tokenData = getTokenData(newTask.Token)
		if len(tokenData[0]) == 0 {
			newTask.Message = "Token data not founded."
		}

		tokenInfo := strings.Split(string(tokenData[1]), "|")
		currTime := getTime()

		newTask.Timestamp = tokenInfo[0]
		newTask.Duration = tokenInfo[1]
		newTask.AccessKey = tokenInfo[2]
		newToken := makeToken(newTask)
		if bytes.Compare(newToken, tokenData[0]) != 0 {
			newTask.Message = "Invalid token."
		} else if !checkDuration(currTime, newTask.Timestamp, parseDuration(newTask.Duration)) {
			newTask.Message = "Time limit excceeded."
		}

		return c.JSON(http.StatusOK, newTask)
	}
}

// loadConfig -
func loadConfig() {
	viper.SetConfigFile("./conf/hmac.yaml")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	task = config{}

	// Reading
	if err := viper.ReadInConfig(); err != nil {
		task = config{}
	}
	// Unmarshal to struct
	if err := viper.Unmarshal(&task); err != nil {
		task = config{}
	}

	secretKey = task.SecretKey
	task.SecretKey = ""
}

// main - Entry point
func main() {
	loadConfig()

	e := echo.New()

	e.File("/", "public/index.html")

	e.GET("/task", getConfigInfo())
	e.PUT("/task", createToken())
	e.GET("/validate", validateToken())

	e.Logger.Fatal(e.Start(":8010"))
}

// ===== [ Public Functions ] =====
