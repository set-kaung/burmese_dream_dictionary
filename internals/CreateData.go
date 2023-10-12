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
	detailMap := make(map[int][]*IndexSearchCache)
	data.Blogs = dream.BlogHeader
	data.SearchData = []string{}
	for _, item := range dream.BlogDetail {
		data.SearchData = append(data.SearchData, item.Content)
		if _, ok := detailMap[item.BlogID]; ok {
			arr := detailMap[item.BlogID]
			detail := &IndexSearchCache{BlogDetailID: item.ID, BlogContent: item.Content}
			arr = append(arr, detail)
			detailMap[item.BlogID] = arr
		} else {
			detail := &IndexSearchCache{BlogDetailID: item.ID, BlogContent: item.Content}
			arr := []*IndexSearchCache{detail}
			detailMap[item.BlogID] = arr
		}
	}
	data.DetailMap = detailMap
	CreateAccentsMap(data)
}

func CreateAccentsMap(data *Data) {
	strs := strings.Split(accents, "\n")
	data.Accents = make(map[string]bool)
	for _, i := range strs {
		data.Accents[i] = true
	}
}
