package xnote

import (
	"github.com/abdulgalimov/xnote/common"
	"github.com/abdulgalimov/xnote/handlers"
	"github.com/abdulgalimov/xnote/web"
)

func Start(xdb common.Db) {
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
