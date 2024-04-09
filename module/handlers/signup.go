package handlers

import (
	"aten/module/models"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"aten/shared/errorx"
	"context"
	"errors"
	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
)

type SignupUser interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
	CreateUser(ctx context.Context, data *models.UserCreate) error
}

type signUpSessionStorage interface {
	SetUserToken(ctx context.Context, userId int, token, subToken string, expiredTime int) error
}

type signupHdl struct {
	store         SignupUser
	sStore        signUpSessionStorage
	tokenProvider tokenprovider.Provider
}

func NewSignupHdl(store SignupUser, sStore signUpSessionStorage, tokenProvider tokenprovider.Provider) *signupHdl {
	return &signupHdl{store: store, sStore: sStore, tokenProvider: tokenProvider}
}

func (h *signupHdl) Response(ctx context.Context, params *models.SignupRequest) (tokenprovider.Token, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"email": params.Email})
	if err == nil && user != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrUserAlreadyExists.Error()).
			WithDebug(errors.New("user already exist").Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	data := models.UserCreate{
		SQLModel: core.NewSQLModel(),
		Email:    params.Email,
		Password: string(hashedPassword),
	}

	if err = h.store.CreateUser(ctx, &data); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   data.Id,
		Email:    data.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := h.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	if err = h.sStore.SetUserToken(ctx, data.Id, accessToken.GetToken(), subToken, expiredTime); err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
