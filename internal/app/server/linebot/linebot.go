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

func Parse(keyword string, items []Response) (*linebotSDK.TemplateMessage, error) {
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
		// char 40 を超えると linebot の仕様上使えないっぽい
		if len(v.Title) > 37 {
			title = v.Title[:37] + "..."
		} else {
			title = v.Title
		}

		resp = append(resp, linebotSDK.NewCarouselColumn(
			v.Image,
			title,
			v.Price,
			linebotSDK.NewURIAction("商品ページ", v.LinkURL),
		))
	}

	return resp
}
