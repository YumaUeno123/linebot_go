package linebot

import (
	"fmt"

	lineBotSDK "github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	Title   string
	Image   string
	LinkURL string
	Price   string
}

const rakutenURL = "https://www.rakuten.co.jp/"
const amazonURL = "https://www.amazon.co.jp/"

func AddSendMessage(kind string, word string, items *[]Response) []lineBotSDK.SendingMessage {
	var sendMessage []lineBotSDK.SendingMessage
	sendMessage = append(sendMessage, lineBotSDK.NewTextMessage(kind+"検索結果"))
	resp, err := parse(kind, word, items)
	if err != nil {
		sendMessage = append(sendMessage, lineBotSDK.NewTextMessage("検索結果がありませんでした"))
	} else {
		sendMessage = append(sendMessage, resp)
	}

	return sendMessage
}

func parse(kind string, keyword string, items *[]Response) (*lineBotSDK.TemplateMessage, error) {
	carouselItems := parseToLineBotFormat(kind, items)
	if len(carouselItems) == 0 {
		return nil, fmt.Errorf("no data")
	}

	resp := lineBotSDK.NewTemplateMessage(
		keyword+"の検索結果",
		lineBotSDK.NewCarouselTemplate(
			carouselItems...,
		),
	)

	return resp, nil
}

func parseToLineBotFormat(kind string, items *[]Response) []*lineBotSDK.CarouselColumn {
	var resp []*lineBotSDK.CarouselColumn

	if items == nil {
		return resp
	}

	for _, v := range *items {
		var title string
		// text の箇所は char 60 を超えると lineBot の仕様上使えないっぽい
		if len(v.Title) > 57 {
			title = v.Title[:57] + "..."
		} else {
			title = v.Title
		}

		// url の長さは char 1000 を超えると lineBot の仕様上使えないっぽい
		var label string
		var url string
		if len(v.LinkURL) > 1000 || len(v.LinkURL) == 0 {
			if kind == "楽天市場" {
				label = "楽天市場で確認"
				url = rakutenURL
			} else {
				label = "Amazonで確認"
				url = amazonURL
			}
		} else {
			label = "商品ページ"
			url = v.LinkURL
		}

		resp = append(resp, lineBotSDK.NewCarouselColumn(
			v.Image,
			v.Price,
			title,
			lineBotSDK.NewURIAction(label, url),
		))
	}

	return resp
}
