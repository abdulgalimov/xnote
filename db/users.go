package db

import (
	"github.com/abdulgalimov/xnote/common"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
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
	ID         int            `db:"id"`
	TelegramID int64          `db:"telegram_id"`
	Name       string         `db:"name"`
	Username   string         `db:"username"`
	Email      string         `db:"email"`
	Password   string         `db:"password"`
	CreatedAt  mysql.NullTime `db:"created_at"`
}

func (u *userInner) user() *common.User {
	return &common.User{
		ID:         u.ID,
		TelegramID: u.TelegramID,
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		Password:   u.Password,
		CreatedAt:  &u.CreatedAt.Time,
	}
}

type dbUsers struct {
	instance *sqlx.DB
}

func (u *dbUsers) Find(userID int) (*common.User, error) {
	var user userInner
	err := u.instance.Get(&user, `SELECT * FROM users WHERE id=? LIMIT 1;`, userID)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) FindByTelegramID(telegramID int64) (*common.User, error) {
	var user userInner
	err := u.instance.Get(&user, `SELECT * FROM users WHERE telegram_id=? LIMIT 1;`, telegramID)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) FindByEmail(email string) (*common.User, error) {
	var user userInner
	err := u.instance.Get(&user, `SELECT * FROM users WHERE email LIKE ? LIMIT 1;`, email)
	if err != nil {
		return nil, err
	}
	return user.user(), nil
}
func (u *dbUsers) Create(src common.User) (*common.User, error) {
	query := `INSERT INTO users (telegram_id, name, username, email, password) VALUES(?, ?, ?, ?, ?);`
	res, err := u.instance.Exec(query, src.TelegramID, src.Name, src.Username, src.Email, src.Password)
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
	user := common.User{
		ID:         int(id),
		TelegramID: src.TelegramID,
		Name:       src.Name,
		Username:   src.Username,
		Email:      src.Email,
		Password:   src.Password,
		CreatedAt:  &now,
	}
	return &user, nil
}
