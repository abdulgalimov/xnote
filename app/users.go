package app

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
	"github.com/abdulgalimov/xnote/core"
	"github.com/abdulgalimov/xnote/models"
)

func userCreate(ctx core.Context) {
	var userSrc models.User
	userSrc.Name = ctx.GetName()
	if ctx.GetTelegramID() == 0 {
		userSrc.Email = ctx.GetEmail()
		userByEmail, _ := xdb.Users().FindByEmail(userSrc.Email)
		if userByEmail != nil {
			ctx.SetError(core.DuplicateError)
			return
		}

		userSrc.Password = ctx.GetPassword()
	} else {
		userSrc.TelegramID = ctx.GetTelegramID()
		userByTelegramID, _ := xdb.Users().FindByTelegramID(userSrc.TelegramID)
		if userByTelegramID != nil {
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
	userID := ctx.GetUserID()
	platform := ctx.GetPlatform()
	deviceID := ctx.GetDeviceID()
	s := fmt.Sprintf("%d|%s|%s|%s|VCcxR6WqmMj2tFnE", userID, platform, deviceID, time.Now().Format(time.RFC3339Nano))
	h := sha1.New()
	h.Write([]byte(s))
	var value = hex.EncodeToString(h.Sum(nil))

	if !isNew {
		oldToken := xdb.Tokens().FindByPlatform(userID, platform, deviceID)
		if oldToken != nil {
			oldToken.Value = value
			xdb.Tokens().Update(oldToken.ID, value)
			return oldToken, nil
		}
	}

	token, err := xdb.Tokens().Create(userID, platform, deviceID, value)
	return token, err
}
