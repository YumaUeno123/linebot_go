package controller

import (
	"net/http"

	"github.com/YumaUeno123/linebot_go/internal/app/service"
)

func Health() {
	http.HandleFunc("/health", service.Health)
}
