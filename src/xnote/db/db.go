package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"xnote/core"
)

var dbInstance *sqlx.DB
type xdb struct {
	users *dbUsers
	tokens *dbTokens
	notes *dbNotes
}
func (d *xdb) Users() core.DbUsers {
	return d.users
}
func (d *xdb) Tokens() core.DbTokens {
	return d.tokens
}
func (d *xdb) Notes() core.DbNotes {
	return d.notes
}

func Connect(config core.DbConnectConfig) (core.Db, error) {
	dbPath := fmt.Sprintf("%s:%d", config.Host, config.Port)
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.UserName, config.Password, dbPath, config.DbName)
	dbTemp, err := sqlx.Connect(config.DriverName, source)
	dbInstance = dbTemp
	//
	initScheme(notesScheme)
	initScheme(usersScheme)
	initScheme(tokensScheme)
	//
	return &xdb{
		users: &dbUsers{},
		tokens: &dbTokens{},
		notes: &dbNotes{},
	}, err
}

func initScheme(scheme string) {
	_, err := dbInstance.Exec(scheme)
	if err != nil {
		panic(err)
	}
}