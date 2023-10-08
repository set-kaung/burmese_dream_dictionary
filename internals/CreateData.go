package internals

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"log"
)

func CreateDataFromJSON(data []byte) *Dream {
	// data, err := os.ReadFile(filename)
	// if err != nil {
	// 	log.Fatalln("Error reading from file with - ", err)
	// }
	// log.Println(string(data[1:20]))
	decoder := json.NewDecoder(bytes.NewReader(data))
	dream := &Dream{}
	if err := decoder.Decode(dream); err != nil {
		log.Fatalln("error decoding JSON:", err)
	}
	return dream
}

func (data *Data) Populate(filedata []byte) {
	dream := CreateDataFromJSON(filedata)
	detailMap := make(map[int][]*DeatailSearchCache)
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
