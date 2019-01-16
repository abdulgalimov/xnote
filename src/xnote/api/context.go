package api

import (
	"encoding/json"
	"net/http"
	"xnote/core"
	"xnote/models"
)

type context struct {
	res   http.ResponseWriter
	req   *http.Request
	reqID string

	cmdCode    core.CmdCode
	userID     int
	text       string
	noteID     int
	name       string
	username   string
	telegramID int64
	email      string
	password   string
	platform   string
	deviceID   string

	pageNum     int
	countOnPage int

	errorCode int

	user       *models.User
	token      *models.Token
	noteList   []*models.Note
	notesCount int
	note       *models.Note

	completeChan chan bool
}

func (ctx *context) GetParam(key string) string {
	return ctx.req.URL.Query().Get(key)
}

func (ctx *context) GetCmdCode() core.CmdCode {
	return ctx.cmdCode
}
func (ctx *context) GetUserID() int {
	return ctx.userID
}
func (ctx *context) GetText() string {
	return ctx.text
}
func (ctx *context) GetName() string {
	return ctx.name
}
func (ctx *context) GetEmail() string {
	return ctx.email
}
func (ctx *context) GetPassword() string {
	return ctx.password
}
func (ctx *context) GetTelegramID() int64 {
	return ctx.telegramID
}
func (ctx *context) GetUsername() string {
	return ctx.username
}
func (ctx *context) GetNoteID() int {
	return ctx.noteID
}
func (ctx *context) GetPageNum() int {
	return ctx.pageNum
}
func (ctx *context) GetCountOnPage() int {
	return ctx.countOnPage
}
func (ctx *context) GetPlatform() string {
	return ctx.platform
}
func (ctx *context) GetDeviceID() string {
	return ctx.deviceID
}

func (ctx *context) SetError(code core.ErrorCode) {
	ctx.errorCode = int(code)
	ctx.Complete()
}
func (ctx *context) SetUser(user *models.User) {
	ctx.user = user
}
func (ctx *context) GetUser() *models.User {
	return ctx.user
}
func (ctx *context) SetNoteList(noteList []*models.Note, count int) {
	ctx.noteList = noteList
	ctx.notesCount = count
}
func (ctx *context) SetNote(note *models.Note) {
	ctx.note = note
}
func (ctx *context) SetToken(token *models.Token) {
	ctx.token = token
}
func (ctx *context) init() {
	ctx.completeChan = make(chan bool)
}
func (ctx *context) Complete() {
	if ctx.completeChan != nil {
		ctx.completeChan <- true
	} else {
		ctx.AnswerData(nil)
	}
}

func (ctx *context) sendResult(res *result) {
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

func (ctx *context) AnswerData(data interface{}) {
	var res result
	if ctx.errorCode == 0 {
		res.Ok = true
		res.Data = data
	} else {
		res.ErrorCode = ctx.errorCode
	}
	ctx.sendResult(&res)
	turns.del(ctx.reqID)

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
	List        []*models.Note `json:"list"`
	Count       int            `json:"count"`
	PageNum     int            `json:"pageNum"`
	CountOnPage int            `json:"countOnPage"`
}
