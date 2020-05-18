package linebot

import (
	"testing"

	"github.com/YumaUeno123/linebot_go/test/matcher"

	lineBotSDK "github.com/line/line-bot-sdk-go/linebot"
)

const (
	keyword = "mock"
	kind    = "mock"
)

var nonItems = []Response{}

var items = []Response{
	{
		Title:   "mock_title",
		Image:   "mock_image",
		Price:   "1,000円",
		LinkURL: "mock_url",
	},
}

func TestParseSendMessage(t *testing.T) {
	type args struct {
		kind    string
		keyword string
		items   []Response
	}
	tests := []struct {
		name string
		args args
		want []lineBotSDK.SendingMessage
	}{
		{
			name: "正常系 items が 0 件場合",
			args: args{
				kind:    kind,
				keyword: keyword,
				items:   nonItems,
			},
			want: []lineBotSDK.SendingMessage{
				lineBotSDK.NewTextMessage(kind + "検索結果"),
				lineBotSDK.NewTextMessage("検索結果がありませんでした"),
			},
		},
		{
			name: "正常系 items が 1 件以上の場合",
			args: args{
				kind:    kind,
				keyword: keyword,
				items:   items,
			},
			want: []lineBotSDK.SendingMessage{
				lineBotSDK.NewTextMessage(kind + "検索結果"),
				lineBotSDK.NewTemplateMessage(
					keyword+"の検索結果",
					lineBotSDK.NewCarouselTemplate(
						lineBotSDK.NewCarouselColumn(
							"mock_image",
							"1,000円",
							"mock_title",
							lineBotSDK.NewURIAction("商品ページ", "mock_url"),
						),
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSendMessage(tt.args.kind, tt.args.keyword, tt.args.items)
			if !matcher.SendMessageResponseMatcher()(got, tt.want) {
				t.Errorf("ParseSendMessage() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
