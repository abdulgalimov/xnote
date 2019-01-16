package app

import (
	"xnote/core"
)

func notesList(ctx core.Context) {
	notesList, count, err := xdb.Notes().FindAll(ctx.GetUserId(), ctx.GetCountOnPage(), ctx.GetPageNum())
	if err != nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetNoteList(notesList, count)
	ctx.Complete()
}

func noteGet(ctx core.Context) {
	noteId := ctx.GetNoteId()
	noteModel, err := xdb.Notes().Find(noteId)
	if err != nil || noteModel == nil {
		ctx.SetError(core.NotFoundError)
		return
	}
	if noteModel.UserId != ctx.GetUserId() {
		ctx.SetError(core.NotFoundError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}
func createNote(ctx core.Context) {
	userId := ctx.GetUserId()
	text := ctx.GetText()
	noteModel, err := xdb.Notes().Create(userId, text)
	if err != nil || noteModel == nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetNote(noteModel)
	ctx.Complete()
}

func deleteNote(ctx core.Context) {
	noteId := ctx.GetNoteId()
	noteModel, err := xdb.Notes().Find(noteId)
	if err != nil || noteModel == nil {
		ctx.SetError(core.NotFoundError)
		return
	}
	if noteModel.UserId != ctx.GetUserId() {
		ctx.SetError(core.NotFoundError)
		return
	}
	err = xdb.Notes().Delete(noteId)
	if err != nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.Complete()
}
