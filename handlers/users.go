package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/abdulgalimov/xnote/common"
	"time"
)

func UserCreate(ctx common.Context, xdb common.Db) {
	var userSrc common.User
	userSrc.Name = ctx.GetName()
	if ctx.GetTelegramID() == 0 {
		userSrc.Email = ctx.GetEmail()
		userByEmail, _ := xdb.Users().FindByEmail(userSrc.Email)
		if userByEmail != nil {
			ctx.SetError(common.DuplicateError)
			return
		}

		userSrc.Password = ctx.GetPassword()
	} else {
		userSrc.TelegramID = ctx.GetTelegramID()
		userByTelegramID, _ := xdb.Users().FindByTelegramID(userSrc.TelegramID)
		if userByTelegramID != nil {
			ctx.SetError(common.DuplicateError)
			return
		}

		userSrc.Username = ctx.GetUsername()
	}

	user, err := xdb.Users().Create(userSrc)
	if err != nil || user == nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetUser(user)

	token, err := createToken(ctx, xdb,true)
	if err != nil || token == nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetToken(token)

	ctx.Complete()
}

func TokenGet(ctx common.Context, xdb common.Db) {
	user := ctx.GetUser()
	if user.Email != ctx.GetEmail() || user.Password != ctx.GetPassword() {
		ctx.SetError(common.AccessError)
		return
	}
	token, err := createToken(ctx, xdb,false)
	if err != nil || token == nil {
		ctx.SetError(common.SystemError)
		return
	}
	ctx.SetToken(token)

	ctx.Complete()
}

func createToken(ctx common.Context, xdb common.Db, isNew bool) (*common.Token, error) {
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
