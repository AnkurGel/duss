package algo

import (
	"crypto/sha1"
	"fmt"
	"github.com/goware/urlx"
	"github.com/lytics/base62"
	"strings"
)

type UrlValidationError struct {
	Reason string
}

func (e *UrlValidationError) Error() string {
	return fmt.Sprintf("UrlValidationError: %s", e.Reason)
}

// This function will iteratively improve
func NormalizeUrl(str string) (string, error) {
	strings.Replace(str, " ", "", -1)
	if len(str) < 4 || len(str) > 2048 {
		return "", &UrlValidationError{"URL has inadequate length"}
	}
	val, err := urlx.Parse(str)
	if err != nil {
		return "", &UrlValidationError{"Cannot parse. Invalid URL."}
	}

	normalize, err := urlx.Normalize(val)
	if err != nil {
		return "", &UrlValidationError{"Cannot normalize URL"}
	}
	return normalize, nil
}

func ComputeHash(str string, offset int) string {
	sha := sha1.Sum([]byte(str))
	hash := base62.StdEncoding.EncodeToString(sha[:])
	startIndex := offset
	endIndex := 7 + startIndex
	return hash[startIndex:endIndex]
}