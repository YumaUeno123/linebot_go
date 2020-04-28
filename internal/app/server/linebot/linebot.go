package linebot

import (
	"fmt"

	linebotSDK "github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	Title   string
	Image   string
	LinkURL string
	Price   string
}

func AddSendMessage(kind string, word string, items []Response) []linebotSDK.SendingMessage {
	var sendMessage []linebotSDK.SendingMessage
	sendMessage = append(sendMessage, linebotSDK.NewTextMessage(kind+"検索結果"))
	resp, err := parse(word, items)
	if err != nil {
		sendMessage = append(sendMessage, linebotSDK.NewTextMessage("検索結果がありませんでした"))
	} else {
		sendMessage = append(sendMessage, resp)
	}

	return sendMessage
}

func parse(keyword string, items []Response) (*linebotSDK.TemplateMessage, error) {
	carouselItems := parseToLinebotFormat(items)
	if len(carouselItems) == 0 {
		return nil, fmt.Errorf("no data")
	}

	resp := linebotSDK.NewTemplateMessage(
		keyword+"の検索結果",
		linebotSDK.NewCarouselTemplate(
			carouselItems...,
		),
	)

	return resp, nil
}

func parseToLinebotFormat(items []Response) []*linebotSDK.CarouselColumn {
	var resp []*linebotSDK.CarouselColumn

	if items == nil {
		return resp
	}

	for _, v := range items {
		var title string
		// text の箇所は char 60 を超えると linebot の仕様上使えないっぽい
		if len(v.Title) > 57 {
			title = v.Title[:57] + "..."
		} else {
			title = v.Title
		}

		resp = append(resp, linebotSDK.NewCarouselColumn(
			v.Image,
			v.Price,
			title,
			linebotSDK.NewURIAction("商品ページ", v.LinkURL),
		))
	}

	return resp
}
