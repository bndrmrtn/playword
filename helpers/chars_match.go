package helpers

import (
	"errors"
	"slices"
	"strings"
)

// CharsMatch Checks if the 2 given strings has the same characters or not.
func CharsMatch(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	chars1 := strings.Split(s1, "")
	chars2 := strings.Split(s2, "")
	succeed := true
	for _, char := range chars1 {
		if slices.Contains(chars2, char) {
			if key, err := findKey(chars2, char); err != nil {
				succeed = false
			} else {
				chars2 = removeKey(chars2, key)
			}
		} else {
			succeed = false
		}
	}
	return succeed
}

// Helpers
func removeKey(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func findKey(slice []string, key string) (int, error) {
	for i, value := range slice {
		if value == key {
			return i, nil
		}
	}
	return 0, errors.New("key not found")
}
