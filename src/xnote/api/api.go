package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"xnote/core"
)

const (
	ErrorInvalidToken = -1000 - iota
	ErrorTurnRequest
	ErrorInvalidUserId
	ErrorInvalidText
	ErrorInvalidNoteId
	ErrorInvalidTelegramId
	ErrorDuplicateTelegramId
	ErrorInvalidName
	ErrorInvalidPlatform
	ErrorInvalidDeviceId
	ErrorInvalidEmail
	ErrorDuplicateEmail
	ErrorInvalidPassword
)

var channel core.ContextReader
var resolver *Resolver

var xdb core.Db
func Init(db core.Db) {
	xdb = db
	channel = make(core.ContextReader, 10)
	resolver = &Resolver{
		handlers: 	make(map[string]HandlerFunc),
		flags: 		make(map[string]int),
		cache:    	make(map[string]*regexp.Regexp),
	}
	resolver.add("^GET /user/create/$", userCreate, flagNoUser | flagNoToken)
	resolver.add("^GET /user/token/$", userTokenGet, flagNoToken)
	resolver.add("^GET /notes/$", notesList, 0)
	resolver.add("^GET /note/$", noteGet, 0)
	resolver.add("^GET /note/create/$", noteCreate, 0)
	resolver.add("^GET /note/delete/$", noteDelete, 0)
}

func GetContextReader() core.ContextReader {
	return channel
}
func Start() {
	err := http.ListenAndServe(":9898", resolver)
	if err != nil {
		panic(err)
	}
}

type HandlerFunc func(ctx *Context)

type Resolver struct {
	handlers map[string]HandlerFunc
	flags	 map[string]int
	cache    map[string]*regexp.Regexp
}

const (
	flagNoUser  = 0x1
	flagNoToken = 0x2
)
func (r *Resolver) add(regex string, handler HandlerFunc, flags int) {
	r.handlers[regex] = handler
	r.flags[regex] = flags
	cache, _ := regexp.Compile(regex)
	r.cache[regex] = cache
}

type UserTurns struct {
	sync.Mutex
	reqs map[string]bool
}
var userTurns = UserTurns{
	reqs: make(map[string]bool),
}
func (u *UserTurns) has(reqId string) bool {
	return u.reqs[reqId]
}
func (u *UserTurns) add(reqId string) {
	u.Lock()
	u.reqs[reqId] = true
	u.Unlock()
}
func (u *UserTurns) del(reqId string) {
	u.Lock()
	u.reqs[reqId] = false
	u.Unlock()
}

func (r *Resolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(check) == true {
			flags := r.flags[pattern]
			var ctx Context
			ctx.res = res
			ctx.req = req

			if flags & flagNoUser == 0 {
				if !loadUser(&ctx) {
					return
				}
				ctx.reqId = fmt.Sprintf("%d-%s", ctx.GetUserId(), check)
				if userTurns.has(ctx.reqId) {
					ctx.SetError(ErrorTurnRequest)
					return
				}
				userTurns.add(ctx.reqId)
			}
			if flags & flagNoToken == 0 {
				if !loadToken(&ctx) {
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

func loadUser(ctx *Context) bool {
	userIdStr := ctx.GetParam("user_id")
	userId, _ := strconv.Atoi(userIdStr)
	if userId <= 0 {
		ctx.SetError(ErrorInvalidUserId)
		return false
	}

	ctx.userId = userId
	//
	user, _ := xdb.Users().Find(ctx.userId)
	if user == nil {
		ctx.SetError(ErrorInvalidUserId)
		return false
	}
	ctx.SetUser(user)
	return true
}

func loadToken(ctx *Context) bool {
	tokenValue := ctx.GetParam("token")
	if tokenValue == "" {
		ctx.SetError(ErrorInvalidToken)
		return false
	}
	token := xdb.Tokens().FindByValue(tokenValue)
	if token == nil || token.UserId != ctx.GetUserId() {
		ctx.SetError(ErrorInvalidToken)
		return false
	}
	ctx.SetToken(token)
	return true
}


func parsePages(ctx *Context) {
	countOnPageStr := ctx.GetParam("countOnPage")
	if countOnPageStr != "" {
		countOnPage, err := strconv.Atoi(countOnPageStr)
		if err != nil {
			return
		}
		pageNumStr := ctx.GetParam("pageNum")
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

func dropContext(ctx *Context) {
	channel <- ctx
}
