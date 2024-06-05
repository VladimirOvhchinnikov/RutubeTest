package telegramconnect

import (
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TelegramClient struct {
	Logger *zap.Logger
	Bot    *tgbotapi.BotAPI
}

func NewTelegramClient(logger *zap.Logger) (*TelegramClient, error) {

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		logger.Error("TELEGRAM_BOT_TOKEN not set in .env file")
		return nil, errors.New("TELEGRAM_BOT_TOKEN not set in .env file")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("Failed to create Telegram bot", zap.Error(err))
		return nil, err
	}

	return &TelegramClient{
		Logger: logger,
		Bot:    bot,
	}, nil
}

func (tc *TelegramClient) GetUserInfo(userID int64) (*tgbotapi.Chat, error) {

	var config tgbotapi.ChatInfoConfig = tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: userID},
	}

	chat, err := tc.Bot.GetChat(config)
	if err != nil {
		tc.Logger.Error("Failed to get chat info from Telegram", zap.Error(err))
		return nil, err
	}
	return &chat, nil
}

func (tc *TelegramClient) Response(userID int64, message string) error {

	msg := tgbotapi.NewMessage(userID, message)
	_, err := tc.Bot.Send(msg)
	if err != nil {
		tc.Logger.Error("Error sending message to user", zap.Error(err))
	}

	return nil
}
