package db

import (
	"github.com/go-sql-driver/mysql"
	"time"
	"xnote/models"
)

var usersScheme = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT,
    name varchar(128) NULL DEFAULT NULL,
	username varchar(128) NULL DEFAULT NULL,
	email varchar(128) NULL DEFAULT NULL,
	password varchar(128) NULL DEFAULT NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
`

type userInner struct {
	Id         int            `db:"id"`
	TelegramId int64          `db:"telegram_id"`
	Name       string         `db:"name"`
	Username   string         `db:"username"`
	Email      string         `db:"email"`
	Password   string         `db:"password"`
	CreatedAt  mysql.NullTime `db:"created_at"`
}

func (u *userInner) user() *models.User {
	return &models.User{
		Id:         u.Id,
		TelegramId: u.TelegramId,
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		Password:   u.Password,
		CreatedAt:  &u.CreatedAt.Time,
	}
}

type dbUsers struct {
}

func (u *dbUsers) Find(userId int) (*models.User, error) {
	var user userInner
	err := dbInstance.Get(&user, `SELECT * FROM users WHERE id=? LIMIT 1;`, userId)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) FindByTelegramId(telegramId int64) (*models.User, error) {
	var user userInner
	err := dbInstance.Get(&user, `SELECT * FROM users WHERE telegram_id=? LIMIT 1;`, telegramId)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) FindByEmail(email string) (*models.User, error) {
	var user userInner
	err := dbInstance.Get(&user, `SELECT * FROM users WHERE email LIKE ? LIMIT 1;`, email)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) Create(src models.User) (*models.User, error) {
	query := `INSERT INTO users (telegram_id, name, username, email, password) VALUES(?, ?, ?, ?, ?);`
	res, err := dbInstance.Exec(query, src.TelegramId, src.Name, src.Username, src.Email, src.Password)
	if err != nil {
		return nil, err
	}
	//
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	//
	now := time.Now()
	user := models.User{
		Id:         int(id),
		TelegramId: src.TelegramId,
		Name:       src.Name,
		Username:   src.Username,
		Email:      src.Email,
		Password:   src.Password,
		CreatedAt:  &now,
	}
	return &user, nil
}
