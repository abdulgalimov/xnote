package web

import (
	"github.com/abdulgalimov/xnote/common"
	"strconv"
)

func userCreate(ctx *context) {
	ctx.name = ctx.GetParam("name")
	if ctx.name == "" {
		ctx.SetError(errorInvalidName)
		return
	}

	ctx.platform = ctx.GetParam("platform")
	if ctx.platform == "" {
		ctx.SetError(errorInvalidPlatform)
		return
	}
	ctx.deviceID = ctx.GetParam("device_id")
	if ctx.deviceID == "" {
		ctx.SetError(errorInvalidDeviceID)
		return
	}

	if ctx.GetParam("telegram_id") == "" {
		ctx.email = ctx.GetParam("email")
		if ctx.email == "" {
			ctx.SetError(errorInvalidEmail)
			return
		}
		ctx.password = ctx.GetParam("password")
		if ctx.password == "" {
			ctx.SetError(errorInvalidPassword)
			return
		}
	} else {
		tgIDStr := ctx.GetParam("telegram_id")
		telegramID, _ := strconv.ParseInt(tgIDStr, 10, 64)
		if telegramID == 0 {
			ctx.SetError(errorInvalidTelegramID)
			return
		}
		ctx.telegramID = telegramID
		ctx.username = ctx.GetParam("username")
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
	ctx.email = ctx.GetParam("email")
	if ctx.email == "" {
		ctx.SetError(errorInvalidEmail)
		return
	}
	ctx.password = ctx.GetParam("password")
	if ctx.password == "" {
		ctx.SetError(errorInvalidPassword)
		return
	}
	ctx.platform = ctx.GetParam("platform")
	if ctx.platform == "" {
		ctx.SetError(errorInvalidPlatform)
		return
	}
	ctx.deviceID = ctx.GetParam("device_id")
	if ctx.deviceID == "" {
		ctx.SetError(errorInvalidDeviceID)
		return
	}

	ctx.cmdCode = common.UserTokenGetCmd
	dropContext(ctx)

	<-ctx.completeChan

	ctx.AnswerData(ctx.token)
}
