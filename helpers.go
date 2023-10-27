package main

import (
	"dream_dictionary/internals"
	"log"
	"strings"
	"sync"
)

const (
	TALL_AA   rune = 'ါ'
	AA             = 'ာ'
	I              = 'ိ'
	II             = 'ီ'
	U              = 'ု'
	UU             = 'ူ'
	E              = 'ေ'
	AI             = 'ဲ'
	ANUSVARA       = 'ံ'
	DOT_BELOW      = '့'
	VISARGA        = 'း'
	VIRAMA         = '္'
	ASAT           = '်'
	MEDIAL_YA      = 'ျ'
	MEDIAL_RA      = 'ြ'
	MEDIAL_WA      = 'ွ'
	MEDIAL_HA      = 'ှ'
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
	sRunes := []rune(source)
	qRunes := []rune(query)
	sWords := map[string][]int{}
	qWords := map[string][]int{}
	//individual words of the source string
	sWords = SplitIntoWords(data.Diacritics_Map, sWords, sRunes)
	//individual words of the query string
	qWords = SplitIntoWords(data.Diacritics_Map, qWords, qRunes)

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

func IsParticle(data *internals.Data, next_char string) bool {
	if data.Accents[next_char] {
		return true
	}
	return false
}

func SplitIntoWords(diacritics_map map[rune]string, words map[string][]int, sRunes []rune) map[string][]int {
	index := 0
	builder := strings.Builder{}
	var nextRune rune
	for i := 0; i < len(sRunes); i++ {
		r := sRunes[i]

		//checking whether index out of bounds
		//for nextRune.
		//if out of bound current rune
		//and next is the same
		//this is current rune is the last one
		if i != len(sRunes)-1 {
			nextRune = sRunes[i+1]
		} else {
			nextRune = r
		}
		if _, ok := diacritics_map[r]; ok {
			//we skipping checking if next rune is diacritic
			//if currnent rune is ္
			if r != VIRAMA {
				if _, ok = diacritics_map[nextRune]; !ok {
					builder.WriteRune(r)
					if i+2 <= len(sRunes)-1 {
						//we check if the next rune is
						//something like တ်
						//if it is, then current word in the buffer is
						//something like နတ်
						n2 := sRunes[i+2]
						if n2 == ASAT || n2 == DOT_BELOW {
							continue
						}
						word := builder.String()
						if arr, ok := words[word]; ok {
							arr = append(arr, index)
							words[word] = arr
						} else {
							words[word] = []int{index}
						}
						index++
						builder.Reset()
						continue
					}
				}
			}
		}
		// if all above procedures isn't executed
		//we can safe to assume that current rune
		//is part of a word
		builder.WriteRune(r)

		//if currnent rune is not ္
		//and the next rune is not a diacritics
		//or if the current rune is the last one
		//we do the following
		if _, ok := diacritics_map[nextRune]; !ok && r != VIRAMA || nextRune == r {
			//again checking for something like နတ်
			if i+2 <= len(sRunes)-1 {
				if sRunes[i+2] == ASAT || sRunes[i+2] == DOT_BELOW {
					continue
				}
			}
			word := builder.String()
			if arr, ok := words[word]; ok {
				arr = append(arr, index)
				words[word] = arr
			} else {
				words[word] = []int{index}
			}
			index++
			builder.Reset()
		}
	}
	return words
}
