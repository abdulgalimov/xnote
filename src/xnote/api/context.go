package api

import (
	"encoding/json"
	"net/http"
	"xnote/core"
	"xnote/models"
)

type Context struct {
	res http.ResponseWriter
	req *http.Request
	reqId string

	cmdCode    core.CmdCode
	userId     int
	text       string
	noteId     int
	name       string
	username   string
	telegramId int64
	email      string
	password   string
	platform   string
	deviceId   string

	pageNum 	int
	countOnPage	int

	errorCode int

	user		*models.User
	token		*models.Token
	noteList  	[]*models.Note
	notesCount	int
	note      	*models.Note

	completeChan chan bool
}

func (ctx *Context) GetParam(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *Context) GetCmdCode() core.CmdCode {
	return ctx.cmdCode
}
func (ctx *Context) GetUserId() int {
	return ctx.userId
}
func (ctx *Context) GetText() string {
	return ctx.text
}
func (ctx *Context) GetName() string {
	return ctx.name
}
func (ctx *Context) GetEmail() string {
	return ctx.email
}
func (ctx *Context) GetPassword() string {
	return ctx.password
}
func (ctx *Context) GetTelegramId() int64 {
	return ctx.telegramId
}
func (ctx *Context) GetUsername() string {
	return ctx.username
}
func (ctx *Context) GetNoteId() int {
	return ctx.noteId
}
func (ctx *Context) GetPageNum() int {
	return ctx.pageNum
}
func (ctx *Context) GetCountOnPage() int {
	return ctx.countOnPage
}
func (ctx *Context) GetPlatform() string {
	return ctx.platform
}
func (ctx *Context) GetDeviceId() string {
	return ctx.deviceId
}


func (ctx *Context) SetError(code core.ErrorCode) {
	ctx.errorCode = int(code)
	ctx.Complete()
}
func (ctx *Context) SetUser(user *models.User) {
	ctx.user = user
}
func (ctx *Context) GetUser() *models.User {
	return ctx.user
}
func (ctx *Context) SetNoteList(noteList []*models.Note, count int) {
	ctx.noteList = noteList
	ctx.notesCount = count
}
func (ctx *Context) SetNote(note *models.Note) {
	ctx.note = note
}
func (ctx *Context) SetToken(token *models.Token) {
	ctx.token = token
}
func (ctx *Context) init() {
	ctx.completeChan = make(chan bool)
}
func (ctx *Context) Complete() {
	if ctx.completeChan != nil {
		ctx.completeChan <- true
	} else {
		ctx.AnswerData(nil)
	}
}

func (ctx *Context) sendResult(res *result) {
	js, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	ctx.res.Header().Set("Content-Type", "application/json")
	_, err = ctx.res.Write(js)
	if err != nil {
		panic(err)
	}
}

func (ctx *Context) AnswerData(data interface{}) {
	var res result
	if ctx.errorCode == 0 {
		res.Ok = true
		res.Data = data
	} else {
		res.ErrorCode = ctx.errorCode
	}
	ctx.sendResult(&res)
	userTurns.del(ctx.reqId)

	if ctx.completeChan != nil {
		close(ctx.completeChan)
	}
}



type result struct {
	Ok        bool        `json:"ok"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

type pagesNoteData struct {
	List 		[]*models.Note	`json:"list"`
	Count 		int 			`json:"count"`
	PageNum		int 			`json:"pageNum"`
	CountOnPage	int 			`json:"countOnPage"`
}