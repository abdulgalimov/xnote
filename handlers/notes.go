package handlers

import (
	"github.com/abdulgalimov/xnote/common"
)

// NotesList загрузить список заметок
func NotesList(ctx common.Context, xdb common.Db) {
	var notesList []*common.Note
	var count int
	var err error
	countOnPage := ctx.GetCountOnPage()
	if countOnPage == 0 {
		notesList, err = xdb.Notes().FindAll(ctx.GetUserID())
		if err == nil {
			count = len(notesList)
		}
	} else {
		notesList, count, err = xdb.Notes().FindPage(ctx.GetUserID(), countOnPage, ctx.GetPageNum())
	}
	if err != nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetNoteList(notesList, count)
	ctx.Complete()
}

// NoteGet получить заметку
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

// CreateNote создать заметку
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

// DeleteNote удалить заметку
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
