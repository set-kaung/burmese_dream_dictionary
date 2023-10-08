package main

import (
	"dream_dictionary/internals"
	"strings"
)

func SearchContent(data *internals.Data, query string) []string {
	found := []string{}
	for _, str := range data.SearchData {
		str = strings.ToLower(str)
		str = strings.ReplaceAll(str, " ", "")
		if strings.Contains(str, query) {
			found = append(found, str)
		}
	}

	return found
}
