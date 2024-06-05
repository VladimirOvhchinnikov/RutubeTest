package usecase

import (
	"fmt"
	"regexp"
	telegramconnect "rutube/infrastructure/TelegramConnect"
	"rutube/infrastructure/database"
	"rutube/models"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type UseCase struct {
	Logger *zap.Logger
	db     *database.Database
	tg     *telegramconnect.TelegramClient
}

func NewUseCase(logger *zap.Logger, db *database.Database, tg *telegramconnect.TelegramClient) *UseCase {
	return &UseCase{
		Logger: logger,
		db:     db,
		tg:     tg,
	}
}

func (uc *UseCase) StartCase(firstName string, lastName string, id int) error {
	var userInfo models.ShortUserInfo

	userInfo, err := uc.db.FindUserByID(id)
	if err != nil {
		uc.Logger.Error("Error while querying the database", zap.Error(err))
		return err
	}

	if userInfo == (models.ShortUserInfo{}) {
		text := fmt.Sprintf("Добро пожаловать в бота, %s %s, который позволит отслеживать дни рождения ваших коллег.\n/help - все доступные команды\n/allUser - все коллеги", firstName, lastName)
		uc.tg.Response(int64(id), text)

		newUser := models.ShortUserInfo{
			IDTG:      id,
			FirstName: firstName,
			LastName:  lastName,
			BirthDate: "",
		}
		err := uc.db.InsertUser(newUser)
		if err != nil {
			uc.Logger.Error("Error while adding new user to the database", zap.Error(err))
			return err
		}

		userInfo = newUser
	}

	if userInfo.BirthDate == "" {
		chat, err := uc.tg.GetUserInfo(int64(id))
		if err != nil {
			uc.Logger.Error("The telegram request was not made correctly", zap.Error(err))
			return err
		}

		birthDate, err := findAndFormatDate(chat.Bio)
		if err != nil {
			text := "Мы не нашли информацию о вашем дне рождении в Вашем Био"
			uc.tg.Response(int64(id), text)
			uc.RequestBirthDate(int64(id))
			uc.Logger.Error("Error parsing text or date was not found", zap.Error(err))
		} else {
			err = uc.db.UpdateUserBirthDate(id, birthDate)

			text := "Мы нашли информацию о вашем дне рождении в Вашем Био, надеемся она верная!"
			uc.tg.Response(int64(id), text)

			if err != nil {
				uc.Logger.Error("Error while updating user birth date in the database", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

var datePatterns = []string{
	`(\d{2})[/-](\d{2})[/-](\d{4})`,
	`(\d{4})[/-](\d{2})[/-](\d{2})`,
	`(\d{2})[/-](\d{2})[/-](\d{2})`,
}

func findAndFormatDate(input string) (string, error) {
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if matches != nil {
			dateStr := matches[0]
			dateFormats := []string{
				"02-01-2006", "02/01/2006",
				"2006-01-02", "2006/01/02",
				"02-01-06", "02/01/06",
			}
			for _, format := range dateFormats {
				date, err := time.Parse(format, dateStr)
				if err == nil {
					return date.Format("2006-01-02"), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no valid date found")
}

func (uc *UseCase) RequestBirthDate(userID int64) error {
	text := "Пожалуйста, введите вашу дату рождения в формате ДД-ММ-ГГГГ."
	return uc.tg.Response(userID, text)
}

func (uc *UseCase) SetBirthday(date string, id int) error {
	result, err := findAndFormatDate(date)
	if err != nil {
		uc.Logger.Error("Invalid date format", zap.Error(err))
		text := "Пожалуйста, введите вашу дату рождения в формате ДД-ММ-ГГГГ."
		uc.tg.Response(int64(id), text)
		return err
	}
	err = uc.db.UpdateUserBirthDate(id, result)
	if err != nil {
		uc.Logger.Error("Error updating birth date", zap.Error(err))
		return err
	}
	text := "Спасибо. Информация о вас внесена в список."
	uc.tg.Response(int64(id), text)
	return nil
}

func (uc *UseCase) SetAllUser(id int) error {
	users, err := uc.db.SetAllUser()
	if err != nil {
		return err
	}

	result := JoinUsers(users)
	uc.tg.Response(int64(id), result)
	return nil
}

func JoinUsers(users []models.ShortUserInfo) string {
	var sb strings.Builder

	for _, user := range users {
		sb.WriteString(fmt.Sprintf("ID: %d, TelegramID: %d, FirstName: %s, LastName: %s, BirthDate: %s\n",
			user.ID, user.IDTG, user.FirstName, user.LastName, user.BirthDate))
	}

	return sb.String()
}

func (uc *UseCase) SetSub(id int, idSub string) error {
	subscriberID, err := strconv.ParseInt(idSub, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	user, err := uc.db.FindUserByID(int(subscriberID))
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	if user.IDTG == 0 {
		return fmt.Errorf("user not found")
	}

	subscribed, err := uc.db.IsSubscribed(int64(id), subscriberID)
	if err != nil {
		return fmt.Errorf("error checking subscription: %w", err)
	}

	if subscribed {
		err := uc.db.UnsubscribeFromBirthday(int64(id), subscriberID)
		if err != nil {
			return fmt.Errorf("error unsubscribing: %w", err)
		}
	} else {
		err := uc.db.SubscribeToBirthday(int64(id), subscriberID)
		if err != nil {
			return fmt.Errorf("error subscribing: %w", err)
		}
	}

	return nil
}
