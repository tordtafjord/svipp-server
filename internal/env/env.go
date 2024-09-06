package env

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func LoadEnv(isProd bool) error {
	var envFile string
	if isProd {
		envFile = ".env"
	} else {
		envFile = ".env.development"
	}
	if err := godotenv.Load(envFile); err != nil {
		return err
	}
	return nil
}

func GetString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func GetInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return intValue
}

func GetFloat(key string, defaultValue float64) float64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}

	return floatValue
}

func GetBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}

	return boolValue
}
