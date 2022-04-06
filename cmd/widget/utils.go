package main

import "strings"

// For each element in a slice, apply the function f
func mapS(sl []string, f func(string) string) []string {
	result := make([]string, len(sl))
	for i, s := range sl {
		result[i] = f(s)
	}
	return result
}

// Filter out all strings where the function does not return true
func filterS(sl []string, f func(string) bool) []string {
	result := make([]string, 0, len(sl))
	for i := range sl {
		if f(sl[i]) {
			result = append(result, sl[i])
		}
	}
	return result
}

// With repeated runes that consists of more than 1 byte,
// the positioning of characters in VT100 is not correct;
// there is a space between each rune.
// However, when placing runes at a given x,y, it appears to work.

// Repeat a rune, n number of times.
// Returns an empty string if memory can not be allocated within append.
func RepeatRune(r rune, n uint) string {
	var sb strings.Builder
	for i := uint(0); i < n; i++ {
		_, err := sb.WriteRune(r)
		if err != nil {
			// In the unlikely event that append inside WriteRune won't work
			return ""
		}
	}
	return sb.String()
}

// Repeat a rune, n number of times
func RepeatRune2(r rune, n uint) (string, error) {
	var sb strings.Builder
	for i := uint(0); i < n; i++ {
		_, err := sb.WriteRune(r)
		if err != nil {
			// In the unlikely event that append inside WriteRune won't work
			return "", err
		}
	}
	return sb.String(), nil
}

// Check if a string is not empty
func nonempty(s string) bool {
	return s != ""
}

// Split a string on any newline: \n, \r or \r\n
// Also removes empty lines and trims away whitespace.
func SplitTrim(s string) []string {
	s = strings.Replace(s, "\r", "\n", -1)
	s = strings.Replace(s, "\r\n", "\n", -1)
	return filterS(mapS(strings.Split(s, "\n"), strings.TrimSpace), nonempty)
}
