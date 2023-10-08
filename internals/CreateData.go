package internals

import (
	"encoding/json"
	"log"
	"os"
)

func CreateDataFromJSON(filename string) *Dream {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error reading from file with - ", err)
	}
	// log.Println(string(data[1:20]))
	dream := &Dream{}
	err = json.Unmarshal(data, dream)
	if err != nil {
		log.Fatalln("error unmarshalling json:", err)
	}
	return dream
}

func (data *Data) Populate(filename string) {
	dream := CreateDataFromJSON(filename)
	detailMap := map[int][]*DeatailSearchCache{}
	data.Blogs = dream.BlogHeader
	for _, item := range dream.BlogDetail {
		if _, ok := detailMap[item.BlogID]; ok {
			arr := detailMap[item.BlogID]
			detail := &DeatailSearchCache{BlogDetailID: item.ID, BlogContent: item.Content}
			arr = append(arr, detail)
			detailMap[item.BlogID] = arr
		} else {
			detail := &DeatailSearchCache{BlogDetailID: item.ID, BlogContent: item.Content}
			arr := []*DeatailSearchCache{detail}
			detailMap[item.BlogID] = arr
		}
	}
	data.DetailMap = detailMap
}
