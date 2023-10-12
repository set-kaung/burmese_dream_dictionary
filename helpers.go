package main

import (
	"dream_dictionary/internals"
	"strings"
)

func SearchContent(data *internals.Data, query string) []string {
	if query == "" {
		return []string{}
	}
	found := []string{}
	exacts := []string{}
	for _, str := range data.SearchData {
		str = strings.ToLower(str)
		str = strings.ReplaceAll(str, " ", "")
		if strings.Contains(str, query) {
			found = append(found, str)
			if IsExact(data, query, str) {
				exacts = append(exacts, str)
			}
		}
	}

	if len(exacts) > 0 {
		return exacts
	}
	return found
}

func IsExact(data *internals.Data, query, source string) bool {
	idx := strings.Index(source, query)
	last_c_start, last_c_end := idx+len(query), idx+len(query)+3
	if last_c_end > len(source) {
		return true
	} else {
		if IsParticle(data, source[last_c_start:last_c_end]) {
			return false
		} else {
			return true
		}
	}
}

func IsParticle(data *internals.Data, next_char string) bool {
	if data.Accents[next_char] {
		return true
	}
	return false
}
