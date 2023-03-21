package api

import "github.com/gofiber/fiber/v2"

func (h handler) ConnectToDB(ctx *fiber.Ctx) error {
	h.channels.RunDBConn <- true
	return ctx.SendStatus(fiber.StatusOK)
}
