package database

import (
	"database/sql"
	"fmt"
	"rutube/models"

	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type Database struct {
	Logger *zap.Logger
	DB     *sql.DB
}

func NewDatabase(logger *zap.Logger, db *sql.DB) *Database {

	return &Database{
		Logger: logger,
		DB:     db,
	}

}

func (db *Database) FindUserByID(userID int) (models.ShortUserInfo, error) {
	query, args, err := squirrel.Select("*").From("users").Where(squirrel.Eq{"telegram_id": userID}).ToSql()
	if err != nil {
		db.Logger.Error("Error building SQL query", zap.Error(err))
		return models.ShortUserInfo{}, err
	}

	var user models.ShortUserInfo

	err = db.DB.QueryRow(query, args...).Scan(&user.ID, &user.IDTG, &user.FirstName, &user.LastName, &user.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ShortUserInfo{}, nil // Пользователь не найден
		}

		db.Logger.Error("Error executing query", zap.Error(err))
		return models.ShortUserInfo{}, err
	}

	return user, nil
}

func (db *Database) UpdateUserBirthDate(telegramID int, newBirthDate string) error {

	query, args, err := squirrel.Update("users").
		Set("birth_date", newBirthDate).
		Where(squirrel.Eq{"telegram_id": telegramID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with telegram_id: %d", telegramID)
	}

	return nil
}

func (db *Database) InsertUser(userInfo models.ShortUserInfo) error {

	query, args, err := squirrel.Insert("users").
		Columns("telegram_id", "first_name", "last_name", "birth_date").
		Values(userInfo.IDTG, userInfo.FirstName, userInfo.LastName, userInfo.BirthDate).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (db *Database) SetAllUser() ([]models.ShortUserInfo, error) {

	query, args, err := squirrel.Select("*").From("users").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []models.ShortUserInfo
	for rows.Next() {
		var user models.ShortUserInfo
		err := rows.Scan(&user.ID, &user.IDTG, &user.FirstName, &user.LastName, &user.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil

}

func (db *Database) SubscribeToBirthday(subscriberID, subscribedToID int64) error {
	query, args, err := squirrel.Insert("subscriptions").
		Columns("subscriber_id", "subscribed_to_id").
		Values(subscriberID, subscribedToID).
		ToSql()
	if err != nil {
		db.Logger.Error("Error building SQL query", zap.Error(err))
		return err
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		db.Logger.Error("Error executing query", zap.Error(err))
		return err
	}

	return nil
}

func (db *Database) UnsubscribeFromBirthday(subscriberID, subscribedToID int64) error {
	query, args, err := squirrel.Delete("subscriptions").
		Where(squirrel.Eq{"subscriber_id": subscriberID, "subscribed_to_id": subscribedToID}).
		ToSql()
	if err != nil {
		db.Logger.Error("Error building SQL query", zap.Error(err))
		return err
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		db.Logger.Error("Error executing query", zap.Error(err))
		return err
	}

	return nil
}

func (db *Database) IsSubscribed(subscriberID, subscribedToID int64) (bool, error) {
	query, args, err := squirrel.Select("COUNT(*)").From("subscriptions").
		Where(squirrel.Eq{"subscriber_id": subscriberID, "subscribed_to_id": subscribedToID}).
		ToSql()
	if err != nil {
		db.Logger.Error("Error building SQL query", zap.Error(err))
		return false, err
	}

	var count int
	err = db.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		db.Logger.Error("Error executing query", zap.Error(err))
		return false, err
	}

	return count > 0, nil
}
