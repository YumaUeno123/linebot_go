package rakuten

type APIResponse struct {
	Items []ResponseItem `json:"items"`
}

type ResponseItem struct {
	Item Item `json:"item"`
}

type Item struct {
	MediumImageUrls []ImageUrl `json:"mediumImageUrls"`
	PointRate       int        `json:"pointRate"`
	ItemName        string     `json:"itemName"`
	ItemPrice       int64      `json:"itemPrice"`
	ItemUrl         string     `json:"itemUrl"`
}

type ImageUrl struct {
	ImageUrl string `json:"imageUrl"`
}
