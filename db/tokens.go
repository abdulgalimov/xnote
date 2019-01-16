package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/go-sql-driver/mysql"
	"github.com/abdulgalimov/xnote/common"
	"time"
)

var tokensScheme = `
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
	platform varchar(16),
	device_id varchar(256),
    value varchar(64),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
`

type tokenInner struct {
	ID        int            `db:"id"`
	UserID    int            `db:"user_id"`
	Platform  string         `db:"platform"`
	DeviceID  string         `db:"device_id"`
	Value     string         `db:"value"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
}

func (t *tokenInner) token() *common.Token {
	return &common.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		Platform:  t.Platform,
		DeviceID:  t.DeviceID,
		Value:     t.Value,
		CreatedAt: &t.CreatedAt.Time,
		UpdatedAt: &t.UpdatedAt.Time,
	}
}

type dbTokens struct {
	instance *sqlx.DB
}

func (t *dbTokens) FindByPlatform(userID int, platform string, deviceID string) *common.Token {
	query := `
SELECT * FROM tokens
WHERE user_id=?
	AND platform LIKE ?
	AND device_id LIKE ?
LIMIT 1;`
	var token tokenInner
	err := t.instance.Get(&token, query, userID, platform, deviceID)
	if err != nil {
		return nil
	}
	return token.token()
}

func (t *dbTokens) Update(tokenID int, value string) {
	t.instance.MustExec("UPDATE tokens SET value=?,updated_at=? WHERE id=?", value, time.Now(), tokenID)
}

func (t *dbTokens) FindByValue(value string) *common.Token {
	var token tokenInner
	err := t.instance.Get(&token, `SELECT * FROM tokens WHERE value LIKE ? LIMIT 1;`, value)
	if err != nil {
		return nil
	}
	return token.token()
}
func (t *dbTokens) Create(userID int, platform string, deviceID string, value string) (*common.Token, error) {
	query := `INSERT INTO tokens (user_id, value, platform, device_id) VALUES(?, ?, ?, ?);`
	res, err := t.instance.Exec(query, userID, value, platform, deviceID)
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
	token := common.Token{
		ID:        int(id),
		UserID:    userID,
		Platform:  platform,
		DeviceID:  deviceID,
		Value:     value,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	return &token, nil
}
