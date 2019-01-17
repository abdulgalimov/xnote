package web

import (
	"github.com/abdulgalimov/xnote/common"
	"strconv"
)

func userCreate(ctx *context) {
	ctx.name = ctx.getParam("name")
	if ctx.name == "" {
		ctx.SetError(errorInvalidName)
		return
	}

	ctx.platform = ctx.getParam("platform")
	if ctx.platform == "" {
		ctx.SetError(errorInvalidPlatform)
		return
	}
	ctx.deviceID = ctx.getParam("device_id")
	if ctx.deviceID == "" {
		ctx.SetError(errorInvalidDeviceID)
		return
	}

	if ctx.getParam("telegram_id") == "" {
		ctx.email = ctx.getParam("email")
		if ctx.email == "" {
			ctx.SetError(errorInvalidEmail)
			return
		}
		ctx.password = ctx.getParam("password")
		if ctx.password == "" {
			ctx.SetError(errorInvalidPassword)
			return
		}
	} else {
		tgIDStr := ctx.getParam("telegram_id")
		telegramID, _ := strconv.ParseInt(tgIDStr, 10, 64)
		if telegramID == 0 {
			ctx.SetError(errorInvalidTelegramID)
			return
		}
		ctx.telegramID = telegramID
		ctx.username = ctx.getParam("username")
	}

	ctx.cmdCode = common.UserCreateCmd
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(struct {
		User  *common.User  `json:"user"`
		Token *common.Token `json:"token"`
	}{
		User:  ctx.user,
		Token: ctx.token,
	})
}

func userTokenGet(ctx *context) {
	ctx.email = ctx.getParam("email")
	if ctx.email == "" {
		ctx.SetError(errorInvalidEmail)
		return
	}
	ctx.password = ctx.getParam("password")
	if ctx.password == "" {
		ctx.SetError(errorInvalidPassword)
		return
	}
	ctx.platform = ctx.getParam("platform")
	if ctx.platform == "" {
		ctx.SetError(errorInvalidPlatform)
		return
	}
	ctx.deviceID = ctx.getParam("device_id")
	if ctx.deviceID == "" {
		ctx.SetError(errorInvalidDeviceID)
		return
	}

	ctx.cmdCode = common.UserTokenGetCmd
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.token)
}
