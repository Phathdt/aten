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

type LoginDbStorage interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
}

type loginHandler struct {
	store LoginDbStorage
}

func NewLoginHandler(store LoginDbStorage) *loginHandler {
	return &loginHandler{store: store}
}

func (h *loginHandler) Response(ctx context.Context, params *models.LoginRequest) (string, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"email": params.Email})
	if err != nil {
		return "", core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	userPass := []byte(params.Password)
	dbPass := []byte(user.Password)

	if err = bcrypt.CompareHashAndPassword(dbPass, userPass); err != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrPasswordNotMatch.Error()).
			WithDebug(err.Error())
	}

	day := time.Hour * 24

	// Create the JWT claims, which includes the user ID and expiry time
	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
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
