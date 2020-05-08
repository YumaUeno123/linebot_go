package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	lineBotSDK "github.com/line/line-bot-sdk-go/linebot"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	"github.com/YumaUeno123/linebot_go/internal/app/client/amazon"
	"github.com/YumaUeno123/linebot_go/internal/app/client/rakuten"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	channelSecret = "LINEBOT_CHANNEL_SECRET"
	accessToken   = "LINEBOT_CHANNEL_ACCESS_TOKEN"
)

func LineBot() {
	bot, err := lineBotSDK.New(
		os.Getenv(channelSecret),
		os.Getenv(accessToken),
	)

	if err != nil {
		log.Print(err)
		return
	}

	rakutenClient := rakuten.New("楽天市場")
	amazonClient := amazon.New("Amazon")

	clients := []client.Client{rakutenClient, amazonClient}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("run request")
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == lineBotSDK.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == lineBotSDK.EventTypeMessage {
				switch message := event.Message.(type) {
				case *lineBotSDK.TextMessage:
					var mutex = &sync.Mutex{}
					var sendMessage []lineBotSDK.SendingMessage

					wg := &sync.WaitGroup{}
					for _, v := range clients {
						wg.Add(1)
						go func(client client.Client) {
							res, err := client.Fetch(message.Text)
							if err != nil {
								fmt.Println(err)
							}
							mutex.Lock()
							sendMessage = append(sendMessage, linebot.AddSendMessage(client.GetKind(), message.Text, res)...)
							mutex.Unlock()
							wg.Done()
						}(v)
					}
					c := make(chan struct{})
					go func() {
						defer close(c)
						wg.Wait()
					}()
					select {
					case <-c:
						fmt.Println("fetch success")
					case <-time.After(5 * time.Second):
						if len(sendMessage) == 0 {
							mutex.Lock()
							sendMessage = append(sendMessage, lineBotSDK.NewTextMessage("問題が発生しました。時間を置いて再度お試しください。"))
							mutex.Unlock()
						}
					}

					if _, err := bot.ReplyMessage(event.ReplyToken, sendMessage...).Do(); err != nil {
						log.Print(err)
					}

				default:
					replyMessage := "検索内容をテキストで入力してください"
					if _, err := bot.ReplyMessage(event.ReplyToken, lineBotSDK.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
}
