package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/YumaUeno123/linebot_go/internal/app/client/rakuten"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"

	linebotSDK "github.com/line/line-bot-sdk-go/linebot"
)

const (
	channelSecret = "LINEBOT_CHANNEL_SECRET"
	accessToken   = "LINEBOT_CHANNEL_ACCESS_TOKEN"
)

func LineBot() {
	bot, err := linebotSDK.New(
		os.Getenv(channelSecret),
		os.Getenv(accessToken),
	)

	if err != nil {
		log.Print(err)
		return
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebotSDK.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebotSDK.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebotSDK.TextMessage:
					// client
					resp := rakuten.Fetch(message.Text)

					// server
					var sendMessage []linebotSDK.SendingMessage
					sendMessage = append(sendMessage, linebotSDK.NewTextMessage("楽天市場検索結果"))
					clientResp, err := linebot.Parse(message.Text, resp)
					if err != nil {
						sendMessage = append(sendMessage, linebotSDK.NewTextMessage("検索結果がありませんでした"))
					} else {
						sendMessage = append(sendMessage, clientResp)
					}

					if _, err := bot.ReplyMessage(event.ReplyToken, sendMessage...).Do(); err != nil {
						log.Print(err)
					}

				default:
					replyMessage := "検索内容をテキストで入力してください"
					if _, err := bot.ReplyMessage(event.ReplyToken, linebotSDK.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
}
