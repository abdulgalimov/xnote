package web

import (
	"fmt"
	"github.com/abdulgalimov/xnote/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

const (
	errorInvalidToken = -1000 - iota
	errorTurnRequest
	errorInvalidUserID
	errorInvalidText
	errorInvalidNoteID
	errorInvalidTelegramID
	errorInvalidName
	errorInvalidPlatform
	errorInvalidDeviceID
	errorInvalidEmail
	errorInvalidPassword
)

var channel common.ContextReader
var resolver *httpResolver

// Init инициализировать пакет
func Init(db common.Db) {
	r := gin.Default()
	r.GET("/ping/:id", pingHandler)
	err := r.Run(":9899")
	if err != nil {
		panic(err)
	}


	channel = make(common.ContextReader, 10)
	resolver = &httpResolver{
		db:       db,
		handlers: make(map[string]handlerFunc),
		flags:    make(map[string]int),
		cache:    make(map[string]*regexp.Regexp),
	}
	resolver.add("^GET /user/create/$", userCreate, flagNoUser|flagNoToken)
	resolver.add("^GET /user/token/$", userTokenGet, flagNoToken)
	resolver.add("^GET /notes/$", notesList, 0)
	resolver.add("^GET /note/$", noteGet, 0)
	resolver.add("^GET /note/create/$", noteCreate, 0)
	resolver.add("^GET /note/delete/$", noteDelete, 0)
}

func pingHandler(c *gin.Context) {
	var token common.Token
	token.Value = c.Param("id")
	c.JSON(200, token)
}

// GetContextReader получить канал куда отправляются контексты
func GetContextReader() common.ContextReader {
	return channel
}

// Start стартануть web-сервер
func Start() {
	err := http.ListenAndServe(":9898", resolver)
	if err != nil {
		panic(err)
	}
}

type handlerFunc func(ctx *context)

type httpResolver struct {
	db       common.Db
	handlers map[string]handlerFunc
	flags    map[string]int
	cache    map[string]*regexp.Regexp
}

const (
	flagNoUser  = 0x1
	flagNoToken = 0x2
)

func (r *httpResolver) add(regex string, handler handlerFunc, flags int) {
	r.handlers[regex] = handler
	r.flags[regex] = flags
	cache, _ := regexp.Compile(regex)
	r.cache[regex] = cache
}

type userTurns struct {
	sync.Mutex
	reqs map[string]bool
}

var turns = userTurns{
	reqs: make(map[string]bool),
}

func (u *userTurns) has(reqID string) bool {
	return u.reqs[reqID]
}
func (u *userTurns) add(reqID string) {
	u.Lock()
	u.reqs[reqID] = true
	u.Unlock()
}
func (u *userTurns) del(reqID string) {
	u.Lock()
	u.reqs[reqID] = false
	u.Unlock()
}

func (r *httpResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(check) == true {
			flags := r.flags[pattern]
			var ctx context
			ctx.res = res
			ctx.req = req

			if flags&flagNoUser == 0 {
				if !r.loadUser(&ctx) {
					return
				}
				ctx.reqID = fmt.Sprintf("%d-%s", ctx.GetUserID(), check)
				if turns.has(ctx.reqID) {
					ctx.SetError(errorTurnRequest)
					return
				}
				turns.add(ctx.reqID)
			}
			if flags&flagNoToken == 0 {
				if !r.loadToken(&ctx) {
					return
				}
			}
			ctx.init()
			handlerFunc(&ctx)
			return
		}
	}
	http.NotFound(res, req)
}

func (r *httpResolver) loadUser(ctx *context) bool {
	userIDStr := ctx.getParam("user_id")
	userID, _ := strconv.Atoi(userIDStr)
	if userID <= 0 {
		ctx.SetError(errorInvalidUserID)
		return false
	}

	ctx.userID = userID
	//
	user, _ := r.db.Users().Find(ctx.userID)
	if user == nil {
		ctx.SetError(errorInvalidUserID)
		return false
	}
	ctx.SetUser(user)
	return true
}

func (r *httpResolver) loadToken(ctx *context) bool {
	tokenValue := ctx.getParam("token")
	if tokenValue == "" {
		ctx.SetError(errorInvalidToken)
		return false
	}
	token := r.db.Tokens().FindByValue(tokenValue)
	if token == nil || token.UserID != ctx.GetUserID() {
		ctx.SetError(errorInvalidToken)
		return false
	}
	ctx.SetToken(token)
	return true
}

func parsePages(ctx *context) {
	countOnPageStr := ctx.getParam("countOnPage")
	if countOnPageStr != "" {
		countOnPage, err := strconv.Atoi(countOnPageStr)
		if err != nil {
			return
		}
		pageNumStr := ctx.getParam("pageNum")
		var pageNum = 0
		if pageNumStr != "" {
			pageNum, err = strconv.Atoi(pageNumStr)
			if err != nil {
				return
			}
		}
		ctx.countOnPage = countOnPage
		ctx.pageNum = pageNum
	}
}

func dropContext(ctx *context) {
	channel <- ctx
}
