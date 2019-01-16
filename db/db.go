package db

import (
	"fmt"
	"github.com/xnoteapp/app/common"
	// драйвер для sqlx
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type xdbHelper struct {
	users  *dbUsers
	tokens *dbTokens
	notes  *dbNotes
}

func (d *xdbHelper) Users() common.DbUsers {
	return d.users
}
func (d *xdbHelper) Tokens() common.DbTokens {
	return d.tokens
}
func (d *xdbHelper) Notes() common.DbNotes {
	return d.notes
}

// Connect подключиться к БД
func Connect(config common.DbConnectConfig) (common.Db, error) {
	dbPath := fmt.Sprintf("%s:%d", config.Host, config.Port)
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.UserName, config.Password, dbPath, config.DbName)
	dbInstance, err := sqlx.Connect(config.DriverName, source)
	if err != nil {
		return nil, err
	}
	//
	initScheme(notesScheme, dbInstance)
	initScheme(usersScheme, dbInstance)
	initScheme(tokensScheme, dbInstance)
	//
	xdb := xdbHelper{
		users:  &dbUsers{},
		tokens: &dbTokens{},
		notes:  &dbNotes{},
	}
	xdb.users.instance = dbInstance
	xdb.tokens.instance = dbInstance
	xdb.notes.instance = dbInstance
	return &xdb, nil
}

func initScheme(scheme string, dbInstance *sqlx.DB) {
	_, err := dbInstance.Exec(scheme)
	if err != nil {
		panic(err)
	}
}
