package amazon

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	urlHost    = "www.amazon.co.jp"
	defaultImg = "https://images-fe.ssl-images-amazon.com/images/G/09/gc/designs/livepreview/amazon_dkblue_noto_email_v2016_jp-main._CB462229751_.png"
)

type amazon struct {
	kind      string
	apiClient client.Api
}

func New(kind string) client.Client {
	return &amazon{
		kind:      kind,
		apiClient: client.NewRetry(),
	}
}

func (a *amazon) GetKind() string {
	return a.kind
}

func (a *amazon) Fetch(keyword string) (*[]linebot.Response, error) {
	u := createURL(keyword)
	res, err := a.apiClient.Get(u)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("http status is " + strconv.Itoa(res.StatusCode))
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("run NewDocumentFromReader error")
		return nil, err
	}
	resp := make([]linebot.Response, 0)

	doc.Find(".s-expand-height > .a-section").Each(func(i int, s *goquery.Selection) {
		if i >= client.MaxCarouselNum {
			return
		}
		image, isImg := s.Find("img").Attr("src")
		if !isImg {
			image = defaultImg
		}
		price := s.Find(".a-price-whole").Text()
		linkURL, isLinkURL := s.Find(".a-link-normal").Attr("href")
		title := s.Find("h2 > a > span").Text()

		_resp := linebot.Response{}
		_resp.Title = title
		if price == "" {
			_resp.Price = "商品ページからご確認ください"
		} else {
			_resp.Price = price + "円"
		}
		if !isLinkURL {
			_resp.LinkURL = ""
		} else {
			_resp.LinkURL = client.UrlScheme + "://" + urlHost + linkURL
		}

		_resp.Image = image

		resp = append(resp, _resp)
	})

	return &resp, nil
}

func createURL(keyword string) string {
	u := &url.URL{}
	u.Scheme = client.UrlScheme
	u.Host = urlHost
	u.Path = "s"
	q := u.Query()
	q.Set("k", keyword)
	u.RawQuery = q.Encode()

	return u.String()
}
