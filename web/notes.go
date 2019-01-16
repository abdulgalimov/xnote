package web

import (
	"github.com/xnoteapp/app/common"
	"strconv"
)

func notesList(ctx *context) {
	ctx.cmdCode = common.NotesListCmd
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

func noteGet(ctx *context) {
	noteIDStr := ctx.GetParam("note_id")
	noteID, _ := strconv.Atoi(noteIDStr)
	if noteID <= 0 {
		ctx.SetError(errorInvalidNoteID)
		return
	}
	ctx.cmdCode = common.GetNoteCmd
	ctx.noteID = noteID
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.note)
}

func noteCreate(ctx *context) {
	text := ctx.GetParam("text")
	if text == "" {
		ctx.SetError(errorInvalidText)
		return
	}
	ctx.cmdCode = common.CreateNoteCmd
	ctx.text = text
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.note)
}

func noteDelete(ctx *context) {
	noteIDStr := ctx.GetParam("note_id")
	noteID, _ := strconv.Atoi(noteIDStr)
	if noteID <= 0 {
		ctx.SetError(errorInvalidNoteID)
		return
	}
	ctx.cmdCode = common.DeleteNoteCmd
	ctx.noteID = noteID
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(nil)
}
