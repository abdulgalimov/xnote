package db

import (
	"github.com/go-sql-driver/mysql"
	"time"
	"xnote/models"
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
	Id        int            `db:"id"`
	UserId    int            `db:"user_id"`
	Platform  string         `db:"platform"`
	DeviceId  string         `db:"device_id"`
	Value     string         `db:"value"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
}

func (t *tokenInner) token() *models.Token {
	return &models.Token{
		Id:        t.Id,
		UserId:    t.UserId,
		Platform:  t.Platform,
		DeviceId:  t.DeviceId,
		Value:     t.Value,
		CreatedAt: &t.CreatedAt.Time,
		UpdatedAt: &t.UpdatedAt.Time,
	}
}

type dbTokens struct {
}

func (t *dbTokens) FindByPlatform(userId int, platform string, deviceId string) *models.Token {
	query := `
SELECT * FROM tokens
WHERE user_id=?
	AND platform LIKE ?
	AND device_id LIKE ?
LIMIT 1;`
	var token tokenInner
	err := dbInstance.Get(&token, query, userId, platform, deviceId)
	if err != nil {
		return nil
	}
	return token.token()
}

func (t *dbTokens) Update(tokenId int, value string) {
	dbInstance.MustExec("UPDATE tokens SET value=?,updated_at=? WHERE id=?", value, time.Now(), tokenId)
}

func (t *dbTokens) FindByValue(value string) *models.Token {
	var token tokenInner
	err := dbInstance.Get(&token, `SELECT * FROM tokens WHERE value LIKE ? LIMIT 1;`, value)
	if err != nil {
		return nil
	}
	return token.token()
}
func (t *dbTokens) Create(userId int, platform string, deviceId string, value string) (*models.Token, error) {
	query := `INSERT INTO tokens (user_id, value, platform, device_id) VALUES(?, ?, ?, ?);`
	res, err := dbInstance.Exec(query, userId, value, platform, deviceId)
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
	token := models.Token{
		Id:        int(id),
		UserId:    userId,
		Platform:  platform,
		DeviceId:  deviceId,
		Value:     value,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	return &token, nil
}
