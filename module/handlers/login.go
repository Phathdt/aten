package handlers

import (
	"aten/module/models"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"aten/shared/errorx"
	"context"
	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
)

type LoginDbStorage interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
}

type LoginSessionStorage interface {
	SetUserToken(ctx context.Context, userId int, token, subToken string, expiredTime int) error
}

type loginHandler struct {
	store         LoginDbStorage
	sStore        LoginSessionStorage
	tokenProvider tokenprovider.Provider
}

func NewLoginHandler(store LoginDbStorage, sStore LoginSessionStorage, tokenProvider tokenprovider.Provider) *loginHandler {
	return &loginHandler{store: store, sStore: sStore, tokenProvider: tokenProvider}
}

func (h *loginHandler) Response(ctx context.Context, params *models.LoginRequest) (tokenprovider.Token, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"email": params.Email})
	if err != nil {
		return nil, core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	userPass := []byte(params.Password)
	dbPass := []byte(user.Password)

	if err = bcrypt.CompareHashAndPassword(dbPass, userPass); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrPasswordNotMatch.Error()).
			WithDebug(err.Error())
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   user.Id,
		Email:    user.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := h.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	if err = h.sStore.SetUserToken(ctx, user.Id, accessToken.GetToken(), subToken, expiredTime); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
