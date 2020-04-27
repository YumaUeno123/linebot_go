package rakuten

type APIResponse struct {
	Items []ResponseItem `json:"items"`
}

type ResponseItem struct {
	Item Item `json:"item"`
}

type Item struct {
	MediumImageUrls    []ImageUrl `json:"mediumImageUrls"`
	PointRate          int        `json:"pointRate"`
	ShopOfTheYearFlag  int        `json:"shopOfTheYearFlag"`
	AffiliateRate      int        `json:"affiliateRate"`
	ShipOverseasFlag   int        `json:"shipOverseasFlag"`
	AsurakuFlag        int        `json:"asurakuFlag"`
	EndTime            string     `json:"endTime"`
	TaxFlag            int        `json:"taxFlag"`
	StartTime          string     `json:"startTime"`
	ItemCaption        string     `json:"itemCaption"`
	TagIds             []int      `json:"tagIds"`
	SmallImageUrls     []ImageUrl `json:"smallImageUrls"`
	AsurakuClosingTime string     `json:"asurakuClosingTime"`
	ImageFlag          int        `json:"imageFlag"`
	Availability       int        `json:"availability"`
	ShopAffiliateUrl   string     `json:"shopAffiliateUrl"`
	ItemCode           string     `json:"itemCode"`
	PostageFlag        int        `json:"postageFlag"`
	ItemName           string     `json:"itemName"`
	ItemPrice          int64      `json:"itemPrice"`
	PointRateEndTime   string     `json:"pointRateEndTime"`
	ShopCode           string     `json:"shopCode"`
	AffiliateUrl       string     `json:"affiliateUrl"`
	GiftFlag           int        `json:"giftFlag"`
	ShopName           string     `json:"shopName"`
	ReviewCount        int        `json:"reviewCount"`
	AsurakuArea        string     `json:"asurakuArea"`
	ShopUrl            string     `json:"shopUrl"`
	CreditCardFlag     int        `json:"creditCardFlag"`
	ReviewAverage      float64    `json:"reviewAverage"`
	ShipOverseasArea   string     `json:"shipOverseasArea"`
	GenreId            string     `json:"genreId"`
	PointRateStartTime string     `json:"pointRateStartTime"`
	ItemUrl            string     `json:"itemUrl"`
}

type ImageUrl struct {
	ImageUrl string `json:"imageUrl"`
}
