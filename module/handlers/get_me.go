package handlers

import (
	"aten/module/models"
	"aten/shared/errorx"
	"context"
	"github.com/phathdt/service-context/core"
)

type GetMeStorage interface {
	GetUserByCondition(ctx context.Context, cond map[string]interface{}) (*models.User, error)
}

type getMeHdl struct {
	store GetMeStorage
}

func NewGetMeHdl(store GetMeStorage) *getMeHdl {
	return &getMeHdl{store}
}

func (h *getMeHdl) Response(ctx context.Context, userId int) (*models.User, error) {
	user, err := h.store.GetUserByCondition(ctx, map[string]interface{}{"id": userId})
	if err != nil {
		return nil, core.ErrNotFound.
			WithError(errorx.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return user, nil
}
