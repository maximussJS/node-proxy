package config

import (
	env "github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	if err := env.Load(); err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
}

func EnvRequiredString(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s is not set. Check your .env file", key)
	}

	return value
}

func EnvOptionalString(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Printf("%s is not set. Used default value %s", key, defaultValue)
		return defaultValue
	}

	return value
}

func EnvOptionalInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		log.Printf("%s is not set. Used default value %d", key, defaultValue)
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		log.Fatalf("%s is not a number. %s", key, err)
	}

	return valueInt
}

func EnvOptionalInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)

	if value == "" {
		log.Printf("%s is not set. Used default value %d", key, defaultValue)
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		log.Fatalf("%s is not a number. %s", key, err)
	}

	return int64(valueInt)
}

func EnvOptionalBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)

	if value == "" {
		log.Printf("%s is not set. Used default value %t", key, defaultValue)
		return defaultValue
	}

	valueBool, err := strconv.ParseBool(value)

	if err != nil {
		log.Fatalf("%s is not a boolean. %s", key, err)
	}

	return valueBool
}
