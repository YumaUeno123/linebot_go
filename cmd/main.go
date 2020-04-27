package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/YumaUeno123/linebot_go/internal/app/controller"
)

var port = flag.String("port", "8080", "http port")

func main() {
	flag.Parse()

	controller.Health()
	controller.LineBot()

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
