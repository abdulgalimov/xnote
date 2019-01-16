package xnote

import (
	"xnote/api"
	"xnote/app"
	"xnote/core"
	"xnote/db"
)

func Start() {
	var dbConfig core.DbConnectConfig
	dbConfig.Host = "localhost"
	dbConfig.Port = 3306
	dbConfig.DriverName = "mysql"
	dbConfig.UserName = "root"
	dbConfig.Password = "123"
	dbConfig.DbName = "xnote_dev"
	xdb, err := db.Connect(dbConfig)
	if err != nil {
		panic(err)
	}
	app.Init(xdb)
	api.Init(xdb)
	go api.Start()
	readContext(api.GetContextReader())
}

func readContext(contextReader core.ContextReader) {
	for {
		ctx := <-contextReader
		app.ParseContext(ctx)
	}
}
