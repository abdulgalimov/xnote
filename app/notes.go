package app

import (
	"github.com/abdulgalimov/xnote/core"
)

func notesList(ctx core.Context) {
	notesList, count, err := xdb.Notes().FindAll(ctx.GetUserID(), ctx.GetCountOnPage(), ctx.GetPageNum())
	if err != nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetNoteList(notesList, count)
	ctx.Complete()
}

func noteGet(ctx core.Context) {
	noteID := ctx.GetNoteID()
	noteModel, err := xdb.Notes().Find(noteID)
	if err != nil || noteModel == nil {
		ctx.SetError(core.NotFoundError)
		return
	}
	if noteModel.UserID != ctx.GetUserID() {
		ctx.SetError(core.NotFoundError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}
func createNote(ctx core.Context) {
	userID := ctx.GetUserID()
	text := ctx.GetText()
	noteModel, err := xdb.Notes().Create(userID, text)
	if err != nil || noteModel == nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}

func deleteNote(ctx core.Context) {
	noteID := ctx.GetNoteID()
	noteModel, err := xdb.Notes().Find(noteID)
	if err != nil || noteModel == nil {
		ctx.SetError(core.NotFoundError)
		return
	}
	if noteModel.UserID != ctx.GetUserID() {
		ctx.SetError(core.NotFoundError)
		return
	}
	err = xdb.Notes().Delete(noteID)
	if err != nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.Complete()
}
