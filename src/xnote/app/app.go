package app

import (
	"xnote/core"
)

var xdb core.Db

// Init инициализировать пакет
func Init(db core.Db) {
	xdb = db
}

// ParseContext обработать полученный запрос
func ParseContext(ctx core.Context) {
	switch ctx.GetCmdCode() {
	case core.NotesListCmd:
		notesList(ctx)
		break
	case core.GetNoteCmd:
		noteGet(ctx)
		break
	case core.CreateNoteCmd:
		createNote(ctx)
		break
	case core.DeleteNoteCmd:
		deleteNote(ctx)
		break

	case core.UserCreateCmd:
		userCreate(ctx)
		break
	case core.UserTokenGetCmd:
		tokenGet(ctx)
		break
	}
}
