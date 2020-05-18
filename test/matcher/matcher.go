package matcher

import (
	"reflect"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Matcher func(got, want interface{}) bool

func SendMessageResponseMatcher() Matcher {
	return func(got, want interface{}) bool {
		if (got == nil && want == nil) || (reflect.ValueOf(got).IsNil() && reflect.ValueOf(want).IsNil()) {
			return true
		}

		g, gok := got.([]linebot.SendingMessage)

		if !gok {
			return false
		}

		w, wok := want.([]linebot.SendingMessage)
		if !wok {
			return false
		}


		for i, gg := range g {
			switch v := gg.(type) {
			case *linebot.TextMessage:
				if v.Text != w[i].(*linebot.TextMessage).Text {

					return false
				}
			case *linebot.TemplateMessage:
				for j, gotItem := range v.Template.(*linebot.CarouselTemplate).Columns {
					wantItem := w[i].(*linebot.TemplateMessage).Template.(*linebot.CarouselTemplate).Columns[j]
					if gotItem.Text != wantItem.Text {
						return false
					}
					if gotItem.Title != wantItem.Title {
						return false
					}
					if gotItem.ThumbnailImageURL != wantItem.ThumbnailImageURL {
						return false
					}
					if gotItem.Actions[0].(*linebot.URIAction).Label != wantItem.Actions[0].(*linebot.URIAction).Label {
						return false
					}
					if gotItem.Actions[0].(*linebot.URIAction).URI != wantItem.Actions[0].(*linebot.URIAction).URI {
						return false
					}
				}
			}
		}

		return true
	}
}
