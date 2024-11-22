package utils

import "fmt"

func ShortedString(s string, length, start, end int) string {
	shortedString := s
	if len(s) > length {
		shortedString = fmt.Sprintf("%s...%s", s[:start], s[len(s)-end:])
	}
	return shortedString
}
