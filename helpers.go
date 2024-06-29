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
	log.Println("Found contains.")
	return found
}

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
		// log.Println("Found exacts.")
		return exacts
	}
	// log.Println("Found contains.")
	return found
}

func IsWordExact(data *internals.Data, query, source string) (bool, string) {
	//individual words of the source string
	sWords, _ := blitter.Splitter(source)
	//individual words of the query string
	qWords, _ := blitter.Splitter(query)

	totalQueryWords := 0
	for _, v := range qWords {
		totalQueryWords += len(v)
	}
	//need to fix this
	qWordIndices := make([]string, totalQueryWords)
	for k, v := range qWords {
		for _, i := range v {
			qWordIndices[i] = k
		}
	}
	appearsInOrder := false
	if len(qWordIndices) == 1 {
		_, appearsInOrder = sWords[qWordIndices[0]]
	} else {
		for i := 1; i < len(qWordIndices); i++ {
			firstArr := sWords[qWordIndices[i-1]]
			secondArr := sWords[qWordIndices[i]]
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
