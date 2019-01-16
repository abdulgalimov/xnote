package handlers

import (
	"github.com/abdulgalimov/xnote/common"
)

func NotesList(ctx common.Context, xdb common.Db) {
	notesList, count, err := xdb.Notes().FindAll(ctx.GetUserID(), ctx.GetCountOnPage(), ctx.GetPageNum())
	if err != nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetNoteList(notesList, count)
	ctx.Complete()
}

func NoteGet(ctx common.Context, xdb common.Db) {
	noteID := ctx.GetNoteID()
	noteModel, err := xdb.Notes().Find(noteID)
	if err != nil || noteModel == nil {
		ctx.SetError(common.NotFoundError)
		return
	}
	if noteModel.UserID != ctx.GetUserID() {
		ctx.SetError(common.NotFoundError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}
func CreateNote(ctx common.Context, xdb common.Db) {
	userID := ctx.GetUserID()
	text := ctx.GetText()
	noteModel, err := xdb.Notes().Create(userID, text)
	if err != nil || noteModel == nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}

func DeleteNote(ctx common.Context, xdb common.Db) {
	noteID := ctx.GetNoteID()
	noteModel, err := xdb.Notes().Find(noteID)
	if err != nil || noteModel == nil {
		ctx.SetError(common.NotFoundError)
		return
	}
	if noteModel.UserID != ctx.GetUserID() {
		ctx.SetError(common.NotFoundError)
		return
	}
	err = xdb.Notes().Delete(noteID)
	if err != nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.Complete()
}
