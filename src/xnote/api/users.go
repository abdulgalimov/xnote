package api

import (
	"strconv"
	"xnote/core"
	"xnote/models"
)

func userCreate(ctx *Context) {
	ctx.name = ctx.GetParam("name")
	if ctx.name == "" {
		ctx.SetError(ErrorInvalidName)
		return
	}

	ctx.platform = ctx.GetParam("platform")
	if ctx.platform == "" {
		ctx.SetError(ErrorInvalidPlatform)
		return
	}
	ctx.deviceId = ctx.GetParam("device_id")
	if ctx.deviceId == "" {
		ctx.SetError(ErrorInvalidDeviceId)
		return
	}

	if ctx.GetParam("telegram_id") == "" {
		ctx.email = ctx.GetParam("email")
		if ctx.email == "" {
			ctx.SetError(ErrorInvalidEmail)
			return
		}
		ctx.password = ctx.GetParam("password")
		if ctx.password == "" {
			ctx.SetError(ErrorInvalidPassword)
			return
		}
	} else {
		tgIdStr := ctx.GetParam("telegram_id")
		telegramId, _ := strconv.ParseInt(tgIdStr, 10, 64)
		if telegramId == 0 {
			ctx.SetError(ErrorInvalidTelegramId)
			return
		}
		ctx.telegramId = telegramId
		ctx.username = ctx.GetParam("username")
	}

	ctx.cmdCode = core.UserCreateCmd
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(struct {
		User  *models.User  `json:"user"`
		Token *models.Token `json:"token"`
	}{
		User:  ctx.user,
		Token: ctx.token,
	})
}

func userTokenGet(ctx *Context) {
	ctx.email = ctx.GetParam("email")
	if ctx.email == "" {
		ctx.SetError(ErrorInvalidEmail)
		return
	}
	ctx.password = ctx.GetParam("password")
	if ctx.password == "" {
		ctx.SetError(ErrorInvalidPassword)
		return
	}
	ctx.platform = ctx.GetParam("platform")
	if ctx.platform == "" {
		ctx.SetError(ErrorInvalidPlatform)
		return
	}
	ctx.deviceId = ctx.GetParam("device_id")
	if ctx.deviceId == "" {
		ctx.SetError(ErrorInvalidDeviceId)
		return
	}

	ctx.cmdCode = core.UserTokenGetCmd
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.token)
}
