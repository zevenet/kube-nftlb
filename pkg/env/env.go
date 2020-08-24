package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// GetString returns a string given the key from env.
func GetString(key string) string {
	env := os.Getenv(key)
	if env == "" {
		message := fmt.Errorf("%s value is empty", key)
		panic(message)
	}
	return env
}

// GetInt returns an int given the key from env.
func GetInt(key string) int {
	env, err := strconv.Atoi(GetString(key))
	if err != nil {
		panic(err)
	}
	return env
}

// GetTime returns a time duration given they key from env.
func GetTime(key string) time.Duration {
	env, err := time.ParseDuration(GetString(key))
	if err != nil {
		panic(err)
	}
	return env
}

func init() {
	// Read .env and panic if it couldn't be read
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}
