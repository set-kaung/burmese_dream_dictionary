package internals

type Data struct {
	Blogs     []*BlogHeader
	DetailMap map[int][]*DeatailSearchCache
}

type Dream struct {
	BlogHeader []*BlogHeader `json:"BlogHeader"`
	BlogDetail []*Detail     `json:"BlogDetail"`
}

type Detail struct {
	ID      int    `json:"BlogDetailId"`
	BlogID  int    `json:"BlogId"`
	Content string `json:"BlogContent"`
}

type BlogHeader struct {
	BlogId    int    `json:"BlogId"`
	BlogTitle string `json:"BlogTitle"`
}

type DeatailSearchCache struct {
	BlogDetailID int    `json:"DetailID"`
	BlogContent  string `json:"Content"`
}
