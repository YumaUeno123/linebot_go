package client

import (
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	UrlScheme      = "https"
	MaxCarouselNum = 10
)

type Client interface {
	GetKind() string
	Fetch(keyword string) *[]linebot.Response
}
