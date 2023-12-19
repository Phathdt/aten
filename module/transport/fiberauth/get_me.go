package fiberauth

import (
	"aten/module/handlers"
	"aten/module/storage"
	"aten/shared/common"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/core"
	"net/http"
)

func GetMe(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId := ctx.Context().UserValue("userId").(int)
		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		sqlStorage := storage.NewSqlStorage(db)
		hdl := handlers.NewGetMeHdl(sqlStorage)

		user, err := hdl.Response(ctx.Context(), userId)
		if err != nil {
			panic(err)
		}

		user.Mask()

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(user))
	}
}
