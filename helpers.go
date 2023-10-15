package main

import (
	"dream_dictionary/internals"
	"log"
	"strings"
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
			if IsWordExact(data, query, str) {
				exacts = append(exacts, str)
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
	found := []string{}
	exacts := []string{}
	for _, str := range data.SearchData {
		str = strings.ToLower(str)
		str = strings.ReplaceAll(str, " ", "")
		if strings.Contains(str, query) {
			found = append(found, str)
			if IsWordExact(data, query, str) {
				exacts = append(exacts, str)
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

func IsWordExact(data *internals.Data, query, source string) bool {
	sRunes := []rune(source)
	words := map[string]bool{}
	words = SplitIntoWords(data.Diacritics_Map, words, sRunes)
	_, ok := words[query]
	return ok
}

func IsParticle(data *internals.Data, next_char string) bool {
	if data.Accents[next_char] {
		return true
	}
	return false
}

func SplitIntoWords(diacritics_map map[rune]string, words map[string]bool, sRunes []rune) map[string]bool {
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
					//we check if the next rune is
					//something like တ်
					if i+2 <= len(sRunes)-1 {
						//if it is, then current word in the buffer is
						//something like နတ်
						n2 := sRunes[i+2]
						if n2 == ASAT || n2 == DOT_BELOW {
							continue
						}
						words[builder.String()] = true
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
			words[builder.String()] = true
			builder.Reset()
		}
	}
	return words
}
