package handlers

import (
	"aten/module/models"
	"aten/plugins/tokenprovider"
	"aten/shared/common"
	"aten/shared/errorx"
	"context"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
)

type SignupUser interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
	CreateUser(ctx context.Context, data *models.UserCreate) error
}

type signupHdl struct {
	store         SignupUser
	tokenProvider tokenprovider.Provider
}

func NewSignupHdl(store SignupUser, tokenProvider tokenprovider.Provider) *signupHdl {
	return &signupHdl{store: store, tokenProvider: tokenProvider}
}

func (h *signupHdl) Response(ctx context.Context, params *models.SignupRequest) (tokenprovider.Token, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"email": params.Email})
	if err != nil && user != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrUserAlreadyExists.Error()).
			WithDebug(err.Error())
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

	payload := common.TokenPayload{
		UserId: data.Id,
	}

	accessToken, err := h.tokenProvider.Generate(&payload, 3600*24*30)
	if err != nil {
		return nil, core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken, nil
}
