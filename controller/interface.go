package controller

import "net/http"

type HandlersInterface interface {
	CommandHandler(w http.ResponseWriter, r *http.Request)
}
