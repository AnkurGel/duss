// Package algo contains the logic for computing hashes and validating URLs
package algo

import (
	"crypto/sha1"
	"fmt"
	"github.com/goware/urlx"
	"github.com/lytics/base62"
	"strings"
)

// URLValidationError represents standard url validation exception
type URLValidationError struct {
	Reason string
}

func (e *URLValidationError) Error() string {
	return fmt.Sprintf("URLValidationError: %s", e.Reason)
}

// NormalizeURL normalizes the URL or throws error if URl is not parsable
func NormalizeURL(str string) (string, error) {
	strings.Replace(str, " ", "", -1)
	if len(str) < 4 || len(str) > 2048 {
		return "", &URLValidationError{"URL has inadequate length"}
	}
	val, err := urlx.Parse(str)
	if err != nil {
		return "", &URLValidationError{"Cannot parse. Invalid URL."}
	}

	normalize, err := urlx.Normalize(val)
	if err != nil {
		return "", &URLValidationError{"Cannot normalize URL"}
	}
	return normalize, nil
}

// ComputeHash computes base62 of provided string and returns 7 bit string
func ComputeHash(str string, offset int) string {
	sha := sha1.Sum([]byte(str))
	hash := base62.StdEncoding.EncodeToString(sha[:])
	startIndex := offset
	endIndex := 7 + startIndex
	return hash[startIndex:endIndex]
}