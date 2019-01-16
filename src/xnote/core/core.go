package core

import (
	"xnote/models"
)

type CmdCode int
const (
	NotesListCmd CmdCode  = 1
	CreateNoteCmd CmdCode = 2
	DeleteNoteCmd CmdCode = 3
	GetNoteCmd CmdCode    = 4
	UserCreateCmd CmdCode  = 5
	UserTokenGetCmd CmdCode  = 6
)

type ErrorCode int
const (
	SystemError ErrorCode = -2000
	NotFoundError ErrorCode = -2001
	AccessError ErrorCode = -2002
	DuplicateError ErrorCode = -2003
)

type Context interface {
	GetCmdCode() CmdCode
	GetUserId() int
	GetText() string
	GetNoteId() int
	GetName() string
	GetEmail() string
	GetPassword() string
	GetTelegramId() int64
	GetUsername() string
	GetPlatform() string
	GetDeviceId() string

	GetPageNum() int
	GetCountOnPage() int

	SetError(code ErrorCode)
	SetUser(user *models.User)
	GetUser() *models.User
	SetNoteList(noteList []*models.Note, count int)
	SetNote(note *models.Note)
	SetToken(token *models.Token)

	Complete()
}
type ContextReader chan Context


type DbConnectConfig struct {
	Host 		string
	Port 		int
	DriverName 	string
	DbName     	string
	UserName   	string
	Password   	string
}

type Db interface {
	Users() DbUsers
	Tokens() DbTokens
	Notes() DbNotes
}

type DbUsers interface {
	Find(userId int) (*models.User, error)
	FindByTelegramId(telegramId int64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(src models.User) (*models.User, error)
}

type DbNotes interface {
	FindAll(userId int, countOnPage int, pageNum int) ([]*models.Note, int, error)
	Find(noteId int) (*models.Note, error)
	Create(userId int, text string) (*models.Note, error)
	Delete(noteId int) error
}

type DbTokens interface {
	FindByPlatform(userId int, platform string, deviceId string) *models.Token
	Update(noteId int, value string)
	FindByValue(value string) *models.Token
	Create(userId int, platform string, deviceId string, value string) (*models.Token, error)
}