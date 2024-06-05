package controller

import (
	"encoding/json"
	"net/http"
	"rutube/models"
	"rutube/usecase"
	"strings"

	"go.uber.org/zap"
)

type Updates struct {
	models.UserInfo
}

type ApiResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Handlers struct {
	Logger  *zap.Logger
	usecase usecase.UseCaseInterface
	Update  Updates
}

func NewHandlers(logger *zap.Logger, usecase usecase.UseCaseInterface) *Handlers {
	return &Handlers{
		Logger:  logger,
		usecase: usecase,
	}
}

func (h *Handlers) CommandHandler(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&h.Update)
	if err != nil {
		h.Logger.Error("Error in decoding ", zap.Error(err))
		h.sendResponse(w, "Error in decoding: "+err.Error(), http.StatusBadRequest)
		return
	}

	if h.Update.Message.Text != "" {
		err := h.messageHandler(w, r)
		h.Update = Updates{}
		if err != nil {
			h.Logger.Error(err.Error())
			h.sendResponse(w, "Wrong Way", http.StatusBadRequest)
			return
		}
	}

	defer r.Body.Close()
}

func (h *Handlers) messageHandler(w http.ResponseWriter, r *http.Request) error {
	messageParts := strings.SplitN(h.Update.Message.Text, " ", 2)
	command := messageParts[0]
	var param string
	if len(messageParts) > 1 {
		param = messageParts[1]
	}

	switch command {
	case "/start":
		{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			h.startHandler()
		}
	case "/allUser":
		{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			h.setAllUser()
		}
	case "/sub":
		{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			h.setSub(param)
		}
	default:
		{
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			h.setMessage(h.Update.Message.Text)
		}
	}
	return nil
}

func (h *Handlers) startHandler() {
	h.usecase.StartCase(h.Update.Message.From.FirstName, h.Update.Message.From.LastName, int(h.Update.Message.From.ID))
}

func (h *Handlers) setMessage(date string) {
	h.usecase.SetBirthday(date, int(h.Update.Message.From.ID))
}

func (h *Handlers) setAllUser() {
	h.usecase.SetAllUser(int(h.Update.Message.From.ID))
}

func (h *Handlers) setSub(sub string) {
	err := h.usecase.SetSub(int(h.Update.Message.From.ID), sub)
	if err != nil {
		h.Logger.Error("Error in setSub handler", zap.Error(err))
	}
}

func (h *Handlers) sendResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ApiResponse{
		Message: message,
		Success: statusCode >= 200 && statusCode < 300,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Logger.Error("Не получилось создать ответ", zap.Error(err))
	}
}
