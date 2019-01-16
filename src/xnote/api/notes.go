package api

import (
	"strconv"
	"xnote/core"
)

func notesList(ctx *Context) {
	ctx.cmdCode = core.NotesListCmd
	ctx.text = ctx.GetParam("find_text")
	parsePages(ctx)
	dropContext(ctx)

	<-ctx.completeChan

	data := pagesNoteData{
		List:        ctx.noteList,
		Count:       ctx.notesCount,
		PageNum:     ctx.pageNum,
		CountOnPage: ctx.countOnPage,
	}
	ctx.AnswerData(&data)

}

func noteGet(ctx *Context) {
	noteIdStr := ctx.GetParam("note_id")
	noteId, _ := strconv.Atoi(noteIdStr)
	if noteId <= 0 {
		ctx.SetError(ErrorInvalidNoteId)
		return
	}
	ctx.cmdCode = core.GetNoteCmd
	ctx.noteId = noteId
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.note)
}

func noteCreate(ctx *Context) {
	text := ctx.GetParam("text")
	if text == "" {
		ctx.SetError(ErrorInvalidText)
		return
	}
	ctx.cmdCode = core.CreateNoteCmd
	ctx.text = text
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.note)
}

func noteDelete(ctx *Context) {
	noteIdStr := ctx.GetParam("note_id")
	noteId, _ := strconv.Atoi(noteIdStr)
	if noteId <= 0 {
		ctx.SetError(ErrorInvalidNoteId)
		return
	}
	ctx.cmdCode = core.DeleteNoteCmd
	ctx.noteId = noteId
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(nil)
}
