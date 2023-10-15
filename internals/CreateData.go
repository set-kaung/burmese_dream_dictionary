package internals

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
	"strings"
)

//go:embed data/accents.txt
var accents string

//go:embed data/DreamDictionary.json
var file_data []byte

//go:embed data/diacritics.txt
var diacritics string

func CreateDataFromJSON(data []byte) *Dream {
	decoder := json.NewDecoder(bytes.NewReader(data))
	dream := &Dream{}
	if err := decoder.Decode(dream); err != nil {
		log.Fatalln("error decoding JSON:", err)
	}
	return dream
}

func (data *Data) Populate() {
	dream := CreateDataFromJSON(file_data)
	detailMap := make(map[int][]string)
	data.Blogs = dream.BlogHeader
	data.SearchData = []string{}
	for _, item := range dream.BlogDetail {
		data.SearchData = append(data.SearchData, item.Content)
		if _, ok := detailMap[item.BlogID]; ok {
			arr := detailMap[item.BlogID]
			arr = append(arr, item.Content)
			detailMap[item.BlogID] = arr
		} else {
			arr := []string{item.Content}
			detailMap[item.BlogID] = arr
		}
	}
	data.DetailMap = detailMap
	CreateAccentsMap(data)
	CreateDiacriticsMap(data)
}

func CreateAccentsMap(data *Data) {
	strs := strings.Split(accents, "\n")
	data.Accents = make(map[string]bool)
	for _, i := range strs {
		data.Accents[i] = true
	}
}

func CreateDiacriticsMap(data *Data) {
	diacritics = strings.ReplaceAll(diacritics, " ", "")
	diacritics_slice := strings.Split(diacritics, "\n")
	// fmt.Println(diacritics_slice)
	diacritics_map := make(map[rune]string)
	for _, str := range diacritics_slice {
		strs := strings.Split(str, "=")
		for _, r := range strs[1] {
			diacritics_map[r] = strs[0]
		}
	}
	data.Diacritics_Map = diacritics_map
}
