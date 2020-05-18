package linebot

import (
	lineBotSDK "github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	Title   string
	Image   string
	LinkURL string
	Price   string
}

const (
	rakutenURL = "https://www.rakuten.co.jp/"
	amazonURL  = "https://www.amazon.co.jp/"
	noData     = "検索結果がありませんでした"
)

func ParseSendMessage(kind string, keyword string, items []Response) []lineBotSDK.SendingMessage {
	var sendMessage []lineBotSDK.SendingMessage
	sendMessage = append(sendMessage, lineBotSDK.NewTextMessage(kind+"検索結果"))
	if len(items) == 0 {
		sendMessage = append(sendMessage, lineBotSDK.NewTextMessage(noData))
		return sendMessage
	}

	parseItems := parse(kind, items)
	resp := lineBotSDK.NewTemplateMessage(
		keyword+"の検索結果",
		lineBotSDK.NewCarouselTemplate(
			parseItems...,
		),
	)
	sendMessage = append(sendMessage, resp)
	return sendMessage
}

func parse(kind string, items []Response) []*lineBotSDK.CarouselColumn {
	var resp []*lineBotSDK.CarouselColumn

	for _, v := range items {
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
