package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	port                 = "8080"
	channelSecret        = "LINEBOT_CHANNEL_SECRET"
	accessToken          = "LINEBOT_CHANNEL_ACCESS_TOKEN"
	rakutenApplicationID = "RAKUTEN_APPLICATION_ID"
	rakutenUrlScheme     = "https"
	rakutenUrlHost       = "app.rakuten.co.jp"
	ichibaBaseUrlPath    = "services/api/IchibaItem/Search/20170706"
	format               = "json"
)

type RakutenAPIResponse struct {
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

func main() {
	fmt.Println("run main.go")

	bot, err := linebot.New(
		os.Getenv(channelSecret),
		os.Getenv(accessToken),
	)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("run callback func")
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					fmt.Println(message.Text)
					u := generateRequestUrl(message.Text)

					resp, err := http.Get(u)
					if err != nil {
						fmt.Println("err is http get")
						log.Print(err)
					}
					defer resp.Body.Close()

					var responseItems RakutenAPIResponse
					if err := json.NewDecoder(resp.Body).Decode(&responseItems); err != nil {
						fmt.Println("err is json decode")
						return
					}

					var sendMessage []linebot.SendingMessage
					searchWord := linebot.NewTextMessage("楽天市場検索結果")
					clientResp := generateCarouselResponse(message.Text, &responseItems)
					sendMessage = append(sendMessage, searchWord)
					sendMessage = append(sendMessage, clientResp)

					if _, err := bot.ReplyMessage(event.ReplyToken, sendMessage...).Do(); err != nil {
						log.Print(err)
					}

				default:
					replyMessage := "検索内容をテキストで入力してください"
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func generateRequestUrl(keyword string) string {
	u := &url.URL{}
	u.Scheme = rakutenUrlScheme
	u.Host = rakutenUrlHost
	u.Path = ichibaBaseUrlPath
	q := u.Query()
	q.Set("format", format)
	q.Set("keyword", keyword)
	q.Set("applicationId", os.Getenv(rakutenApplicationID))
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	return u.String()
}

func generateCarouselResponse(keyword string, items *RakutenAPIResponse) *linebot.TemplateMessage {
	fmt.Println(items)
	resp := linebot.NewTemplateMessage(
		keyword+"の検索結果",
		linebot.NewCarouselTemplate(
			createCarousel(items)...,
		),
	)

	return resp
}

func createCarousel(items *RakutenAPIResponse) []*linebot.CarouselColumn {
	var resp []*linebot.CarouselColumn

	for index, item := range items.Items {
		if index >= 10 {
			break
		}

		var title string

		if len(item.Item.ItemName) > 40 {
			title = item.Item.ItemName[:40] + "..."
		} else {
			title = item.Item.ItemName
		}

		resp = append(resp, linebot.NewCarouselColumn(
			item.Item.MediumImageUrls[0].ImageUrl,
			title,
			humanize.Comma(item.Item.ItemPrice)+"円",
			linebot.NewURIAction("商品ページ", item.Item.ItemUrl),
			linebot.NewURIAction("ショップサイト", item.Item.ShopUrl),
		))
	}

	return resp
}
