package xnote

import (
	"fmt"
	"github.com/abdulgalimov/xnote/common"
	"github.com/abdulgalimov/xnote/db"
	"github.com/abdulgalimov/xnote/handlers"
	"github.com/abdulgalimov/xnote/web"
)

func Start() {
	fmt.Println("run xnote")
	var dbConfig common.DbConnectConfig
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

	web.Init(xdb)
	go web.Start()
	readContext(web.GetContextReader(), xdb)
}

// ReadContext прочитать запросы
func readContext(contextReader common.ContextReader, xdb common.Db) {
	for {
		ctx := <-contextReader
		parseContext(ctx, xdb)
	}
}


func parseContext(ctx common.Context, xdb common.Db) {
	switch ctx.GetCmdCode() {
	case common.NotesListCmd:
		handlers.NotesList(ctx, xdb)
		break
	case common.GetNoteCmd:
		handlers.NoteGet(ctx, xdb)
		break
	case common.CreateNoteCmd:
		handlers.CreateNote(ctx, xdb)
		break
	case common.DeleteNoteCmd:
		handlers.DeleteNote(ctx, xdb)
		break

	case common.UserCreateCmd:
		handlers.UserCreate(ctx, xdb)
		break
	case common.UserTokenGetCmd:
		handlers.TokenGet(ctx, xdb)
		break
	}
}