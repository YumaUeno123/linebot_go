package rakuten

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dustin/go-humanize"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	model "github.com/YumaUeno123/linebot_go/internal/app/model/rakuten"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	applicationID = "RAKUTEN_APPLICATION_ID"
	urlHost       = "app.rakuten.co.jp"
	urlPath       = "services/api/IchibaItem/Search/20170706"
	format        = "json"
)

type rakuten struct {
	kind string
}

func New(kind string) client.Client {
	return &rakuten{
		kind: kind,
	}
}

func (r *rakuten) GetKind() string {
	return r.kind
}

func (r *rakuten) Fetch(keyword string) []linebot.Response {
	u := createURL(keyword)

	getResp, err := http.Get(u)
	if err != nil {
		log.Print(err)
	}
	defer getResp.Body.Close()

	var responseItems model.APIResponse

	if err := json.NewDecoder(getResp.Body).Decode(&responseItems); err != nil {
		log.Print(err)
	}

	resp := make([]linebot.Response, 0)
	if len(responseItems.Items) == 0 {
		return resp
	}

	var limit int

	if client.MaxCarouselNum > len(responseItems.Items) {
		limit = len(responseItems.Items)
	} else {
		limit = client.MaxCarouselNum
	}

	for i := 0; i < limit; i++ {
		resp = append(resp, parse(&responseItems.Items[i]))
	}

	return resp
}

func parse(responseItem *model.ResponseItem) (resp linebot.Response) {
	resp.Title = responseItem.Item.ItemName
	resp.Image = responseItem.Item.MediumImageUrls[0].ImageUrl
	resp.Price = humanize.Comma(responseItem.Item.ItemPrice) + "円"
	resp.LinkURL = responseItem.Item.ItemUrl
	return
}

func createURL(keyword string) string {
	u := &url.URL{}
	u.Scheme = client.UrlScheme
	u.Host = urlHost
	u.Path = urlPath
	q := u.Query()
	q.Set("format", format)
	q.Set("keyword", keyword)
	q.Set("applicationId", os.Getenv(applicationID))
	u.RawQuery = q.Encode()

	return u.String()
}
