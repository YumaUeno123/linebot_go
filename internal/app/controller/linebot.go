package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/YumaUeno123/linebot_go/internal/app/client/amazon"

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
		fmt.Println("run request")
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
					rakutenChannel := make(chan []linebot.Response)
					amazonChannel := make(chan []linebot.Response)
					go rakuten.Fetch(rakutenChannel, message.Text)
					go amazon.Fetch(amazonChannel, message.Text)
					rakutenResp := <-rakutenChannel
					amazonResp := <-amazonChannel

					var sendMessage []linebotSDK.SendingMessage
					sendMessage = append(sendMessage, linebot.AddSendMessage("楽天市場", message.Text, rakutenResp)...)
					sendMessage = append(sendMessage, linebot.AddSendMessage("amazon", message.Text, amazonResp)...)
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
