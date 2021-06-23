package utils

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ApiPrefix = "https://api2.autopilothq.com/v1"
)

var (
	// ErrorInvalidKeyFormat ...
	ErrorInvalidKeyFormat = errors.New("invalid key format")
)

// GetURLFromParams gets URL from the input key string
func GetURLFromParams(entity, id string) string {
	return fmt.Sprintf("%s/%s/%s", ApiPrefix, entity, id)
}

// GetKeyFromParams get the key used in cache from URL params
func GetKeyFromParams(entity, id string) string {
	return fmt.Sprintf("%s:%s", entity, id)
}

// GetURLFromKey get URL string by key
func GetURLFromKey(key string) (string, error) {
	fields := strings.Split(key, ":")
	if len(fields) < 2 {
		return "", ErrorInvalidKeyFormat
	}

	url := GetURLFromParams(fields[0], fields[1])
	return url, nil
}

// GetPostURLFromKey get POST URL from key string
func GetPostURLFromKey(key string) (string, error) {
	fields := strings.Split(key, ":")
	if len(fields) < 2 {
		return "", ErrorInvalidKeyFormat
	}

	return fmt.Sprintf("%s/%s", ApiPrefix, fields[0]), nil

}
