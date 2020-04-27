package rakuten

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/YumaUeno123/linebot_go/internal/app/model/rakuten"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	rakutenApplicationID = "RAKUTEN_APPLICATION_ID"
	urlScheme            = "https"
	rakutenUrlHost       = "app.rakuten.co.jp"
	ichibaBaseUrlPath    = "services/api/IchibaItem/Search/20170706"
	format               = "json"
	MaxCarouselNum       = 10
)

func Fetch(keyword string) []linebot.Response {
	url := createURL(keyword)

	getResp, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	defer getResp.Body.Close()

	var responseItems rakuten.APIResponse

	if err := json.NewDecoder(getResp.Body).Decode(&responseItems); err != nil {
		log.Print(err)
	}

	if len(responseItems.Items) == 0 {
		return nil
	}

	resp := make([]linebot.Response, 0)
	for i := 0; i < MaxCarouselNum; i++ {
		resp = append(resp, parse(responseItems.Items[i]))
	}

	return resp
}

func parse(responseItem rakuten.ResponseItem) (resp linebot.Response) {
	resp.Title = responseItem.Item.ItemName
	resp.Image = responseItem.Item.MediumImageUrls[0].ImageUrl
	resp.Price = humanize.Comma(responseItem.Item.ItemPrice) + "円"
	resp.LinkURL = responseItem.Item.ItemUrl
	return
}

func createURL(keyword string) string {
	u := &url.URL{}
	u.Scheme = urlScheme
	u.Host = rakutenUrlHost
	u.Path = ichibaBaseUrlPath
	q := u.Query()
	q.Set("format", format)
	q.Set("keyword", keyword)
	q.Set("applicationId", os.Getenv(rakutenApplicationID))
	u.RawQuery = q.Encode()

	return u.String()
}
