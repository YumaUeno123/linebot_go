package amazon

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	urlHost = "www.amazon.co.jp"
)

type amazon struct {
	kind string
}

func New(kind string) client.Client {
	return &amazon{
		kind: kind,
	}
}

func (a *amazon) GetKind() string {
	return a.kind
}

func (a *amazon) Fetch(keyword string) *[]linebot.Response {
	u := createURL(keyword)

	res, err := http.Get(u)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil
	}

	buf, _ := ioutil.ReadAll(res.Body)

	// 文字化け対応
	det := chardet.NewTextDetector()
	detResult, _ := det.DetectBest(buf)

	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detResult.Charset, bReader)

	doc, _ := goquery.NewDocumentFromReader(reader)
	resp := make([]linebot.Response, 0)
	doc.Find(".s-expand-height > .a-section").Each(func(i int, s *goquery.Selection) {
		if i >= client.MaxCarouselNum {
			return
		}
		image, _ := s.Find("img").Attr("src")
		price := s.Find(".a-price-whole").Text()
		linkURL, _ := s.Find(".a-link-normal").Attr("href")
		title := s.Find("h2 > a > span").Text()

		_resp := linebot.Response{}
		_resp.Title = title
		if price == "" {
			_resp.Price = "商品ページからご確認ください"
		} else {
			_resp.Price = price + "円"
		}
		_resp.LinkURL = client.UrlScheme + "://" + urlHost + linkURL
		_resp.Image = image

		resp = append(resp, _resp)
	})

	return &resp
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
