package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
	"xnote/core"
	"xnote/models"
)

func userCreate(ctx core.Context) {
	var userSrc models.User
	userSrc.Name = ctx.GetName()
	if ctx.GetTelegramId() == 0 {
		userSrc.Email = ctx.GetEmail()
		userByEmail, _ := xdb.Users().FindByEmail(userSrc.Email)
		if userByEmail != nil {
			ctx.SetError(core.DuplicateError)
			return
		}

		userSrc.Password = ctx.GetPassword()
	} else {
		userSrc.TelegramId = ctx.GetTelegramId()
		userByTelegramId, _ := xdb.Users().FindByTelegramId(userSrc.TelegramId)
		if userByTelegramId != nil {
			ctx.SetError(core.DuplicateError)
			return
		}

		userSrc.Username = ctx.GetUsername()
	}

	user, err := xdb.Users().Create(userSrc)
	if err != nil || user == nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetUser(user)

	token, err := createToken(ctx, true)
	if err != nil || token == nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetToken(token)

	ctx.Complete()
}

func tokenGet(ctx core.Context) {
	user := ctx.GetUser()
	if user.Email != ctx.GetEmail() || user.Password != ctx.GetPassword() {
		ctx.SetError(core.AccessError)
		return
	}
	token, err := createToken(ctx, false)
	if err != nil || token == nil {
		ctx.SetError(core.SystemError)
		return
	}
	ctx.SetToken(token)

	ctx.Complete()
}

func createToken(ctx core.Context, isNew bool) (*models.Token, error) {
	userId := ctx.GetUserId()
	platform := ctx.GetPlatform()
	deviceId := ctx.GetDeviceId()
	s := fmt.Sprintf("%d|%s|%s|%s|VCcxR6WqmMj2tFnE", userId, platform, deviceId, time.Now().Format(time.RFC3339Nano))
	h := sha1.New()
	h.Write([]byte(s))
	var value = hex.EncodeToString(h.Sum(nil))

	if !isNew {
		oldToken := xdb.Tokens().FindByPlatform(userId, platform, deviceId)
		if oldToken != nil {
			oldToken.Value = value
			xdb.Tokens().Update(oldToken.Id, value)
			return oldToken, nil
		}
	}

	token, err := xdb.Tokens().Create(userId, platform, deviceId, value)
	return token, err
}
