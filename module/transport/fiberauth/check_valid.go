package fiberauth

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"net/http"
)

func CheckValid(sc sctx.ServiceContext) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(core.SimpleSuccessResponse("ok"))
	}
}
