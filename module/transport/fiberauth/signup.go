package fiberauth

import (
	"aten/module/handlers"
	"aten/module/models"
	"aten/module/storage"
	"aten/plugins/tokenprovider"
	"aten/plugins/validation"
	"aten/shared/common"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/gormc"
	"github.com/phathdt/service-context/core"
	"net/http"
)

func SignUp(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var p models.SignupRequest

		if err := ctx.BodyParser(&p); err != nil {
			panic(err)
		}

		if err := validation.Validate(p); err != nil {
			panic(err)
		}

		db := sc.MustGet(common.KeyCompGorm).(gormc.GormComponent).GetDB()
		sqlStorage := storage.NewSqlStorage(db)
		tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)
		hdl := handlers.NewSignupHdl(sqlStorage, tokenProvider)

		token, err := hdl.Response(ctx.Context(), &p)
		if err != nil {
			panic(err)
		}

		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(token))
	}
}
