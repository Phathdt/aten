package handlers

import (
	"aten/module/models"
	"aten/shared/errorx"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type SignupUser interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
	CreateUser(ctx context.Context, data *models.UserCreate) error
}

type signupHdl struct {
	store SignupUser
}

func NewSignupHdl(store SignupUser) *signupHdl {
	return &signupHdl{store: store}
}

func (h *signupHdl) Response(ctx context.Context, params *models.SignupRequest) (string, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"email": params.Email})
	if err != nil && user != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrUserAlreadyExists.Error()).
			WithDebug(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	data := models.UserCreate{
		SQLModel: core.NewSQLModel(),
		Email:    params.Email,
		Password: string(hashedPassword),
	}

	if err = h.store.CreateUser(ctx, &data); err != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrCreateUser.Error()).
			WithDebug(err.Error())
	}

	day := time.Hour * 24

	// Create the JWT claims, which includes the user ID and expiry time
	claims := jwt.MapClaims{
		"id":    data.Id,
		"email": data.Email,
		"exp":   time.Now().Add(day * 1).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", core.ErrInternalServerError.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return t, nil
}
