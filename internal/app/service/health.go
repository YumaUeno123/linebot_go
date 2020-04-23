package service

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("Hello")
}
