package main

import (
	"dream_dictionary/internals"
	"log"
	"strings"
	"sync"

	"github.com/Set-Kaung/blitter"
)

func SearchBlogContents(data *internals.Data, searchStrings []string, query string) []string {
	if query == "" {
		return []string{}
	}

	found := []string{}
	exacts := []string{}
	for _, str := range searchStrings {
		str = strings.ToLower(str)
		str = strings.ReplaceAll(str, " ", "")
		if strings.Contains(str, query) {
			found = append(found, str)
			ok, result := IsWordExact(data, query, str)
			if ok {
				exacts = append(exacts, result)
			}
		}
	}

	if len(exacts) > 0 {
		log.Println("Found exacts.")
		return exacts
	}
	// if can't find exact match, return closet matches
	log.Println("Found contains.")
	return found
}

// need to update this with a worker pool
func SearchContent(data *internals.Data, query string) []string {
	if query == "" {
		return []string{}
	}
	wg := &sync.WaitGroup{}
	found := []string{}
	exacts := []string{}
	for _, str := range data.SearchData {
		str = strings.ToLower(str)
		str = strings.ReplaceAll(str, " ", "")
		if strings.Contains(str, query) {
			found = append(found, str)
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				ok, result := IsWordExact(data, query, s)
				if ok {
					exacts = append(exacts, result)
				}
			}(str)

		}
	}
	wg.Wait()

	if len(exacts) > 0 {
		return exacts
	}
	return found
}

// code to find the exact appearance of the target word
// we are searching for
func IsWordExact(data *internals.Data, query, source string) (bool, string) {
	//individual words of the source string
	sWords, _ := blitter.Splitter(source)
	//individual words of the query string
	qWords, totalQueryWords := blitter.Splitter(query)

	//need to fix this
	queryWordIndices := make([]string, totalQueryWords)
	for k, v := range qWords {
		for _, i := range v {
			queryWordIndices[i] = k
		}
	}

	appearsInOrder := false
	if len(queryWordIndices) == 1 {
		_, appearsInOrder = sWords[queryWordIndices[0]]
	} else {
		for i := 1; i < len(queryWordIndices); i++ {
			firstArr := sWords[queryWordIndices[i-1]]
			secondArr := sWords[queryWordIndices[i]]
			for _, f := range firstArr {
				for _, s := range secondArr {
					if s-f == 1 {
						appearsInOrder = true
						break
					}
				}
			}
			if !appearsInOrder {
				break
			}
		}
	}
	return appearsInOrder, source
}
