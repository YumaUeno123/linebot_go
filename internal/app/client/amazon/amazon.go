package amazon

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

const (
	urlScheme      = "https"
	amazonUrlHost  = "www.amazon.co.jp"
	MaxCarouselNum = 10
)

func Fetch(ch chan<- []linebot.Response, keyword string) {
	url := createURL(keyword)

	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	buf, _ := ioutil.ReadAll(res.Body)

	// 文字化け対応
	det := chardet.NewTextDetector()
	detResult, _ := det.DetectBest(buf)

	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detResult.Charset, bReader)

	doc, _ := goquery.NewDocumentFromReader(reader)
	resp := make([]linebot.Response, 0)
	doc.Find(".s-expand-height").Each(func(i int, s *goquery.Selection) {
		if i >= MaxCarouselNum {
			return
		}
		image, _ := s.Find("img").Attr("src")
		price := s.Find(".a-price-whole").Text()
		linkURL, _ := s.Find(".a-link-normal").Attr("href")
		title := s.Find("h2 > a > span").Text()

		_resp := linebot.Response{}
		_resp.Title = title
		_resp.Price = price + "円"
		_resp.LinkURL = urlScheme + "://" + amazonUrlHost + linkURL
		_resp.Image = image

		resp = append(resp, _resp)
	})

	ch <- resp
}

func createURL(keyword string) string {
	u := &url.URL{}
	u.Scheme = urlScheme
	u.Host = amazonUrlHost
	u.Path = "s"
	q := u.Query()
	q.Set("k", keyword)
	u.RawQuery = q.Encode()

	return u.String()
}
