package core

import (
	"xnote/models"
)

// CmdCode код события
type CmdCode int

// Коды событий приложения
const (
	NotesListCmd CmdCode = iota + 1
	CreateNoteCmd
	DeleteNoteCmd
	GetNoteCmd
	UserCreateCmd
	UserTokenGetCmd
)

// ErrorCode код ошибки приложения
type ErrorCode int

// коды ошибок приложения
const (
	SystemError ErrorCode = iota - 2000
	NotFoundError
	AccessError
	DuplicateError
)

// Context интверфейс, контекст запроса
type Context interface {
	GetCmdCode() CmdCode
	GetUserID() int
	GetText() string
	GetNoteID() int
	GetName() string
	GetEmail() string
	GetPassword() string
	GetTelegramID() int64
	GetUsername() string
	GetPlatform() string
	GetDeviceID() string

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

// ContextReader канал, для отправки запросов
type ContextReader chan Context

// DbConnectConfig конфиг для подключения к БД
type DbConnectConfig struct {
	Host       string
	Port       int
	DriverName string
	DbName     string
	UserName   string
	Password   string
}

// Db интерфейс для работы с БД
type Db interface {
	Users() DbUsers
	Tokens() DbTokens
	Notes() DbNotes
}

// DbUsers интерфейс БД для работы с таблицей users
type DbUsers interface {
	Find(userID int) (*models.User, error)
	FindByTelegramID(telegramID int64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(src models.User) (*models.User, error)
}

// DbNotes интерфейс БД для работы с таблицей notes
type DbNotes interface {
	FindAll(userID int, countOnPage int, pageNum int) ([]*models.Note, int, error)
	Find(noteID int) (*models.Note, error)
	Create(userID int, text string) (*models.Note, error)
	Delete(noteID int) error
}

// DbTokens интерфейс БД для работы с таблицей tokens
type DbTokens interface {
	FindByPlatform(userID int, platform string, deviceID string) *models.Token
	Update(noteID int, value string)
	FindByValue(value string) *models.Token
	Create(userID int, platform string, deviceID string, value string) (*models.Token, error)
}
