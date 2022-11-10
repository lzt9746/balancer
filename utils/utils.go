package utils

import (
	"regexp"
)

// SplitStringBySpaces splits a string by one or more space.
func SplitStringBySpaces(s string) []string {
	r := regexp.MustCompile("[^\\s]+")
	res := r.FindAllString(s, -1)
	return res
}
