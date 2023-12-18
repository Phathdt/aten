package fiberauth

import (
	"aten/module/handlers"
	"aten/module/models"
	"aten/module/storage"
	"aten/plugins/validation"
	"aten/shared/common"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/core"
	"net/http"
)

func Login(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p models.LoginRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		sqlStorage := storage.NewSqlStorage(db)
		hdl := handlers.NewLoginHandler(sqlStorage)

		token, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(map[string]interface{}{
			"token": token,
		}))
	}
}
