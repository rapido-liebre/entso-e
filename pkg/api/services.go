package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h handler) RunServices(ctx *fiber.Ctx) error {
	//check channel availability, open it if it's closed
	//if h.ChannelIsClosed(h.channels.ArchIsRunning) {
	//	h.channels.ArchIsRunning = make(chan bool)
	//}
	//if h.ChannelIsClosed(h.channels.WatchIsRunning) {
	//	h.channels.WatchIsRunning = make(chan bool)
	//}
	// send confirmation all services are running up
	h.channels.ProcessorIsRunning <- true
	h.channels.ParserIsRunning <- true

	h.channels.RunParse <- true

	return ctx.SendStatus(fiber.StatusOK)
}

func (h handler) StopServices(ctx *fiber.Ctx) error {
	h.channels.ParserIsRunning <- false
	h.channels.ProcessorIsRunning <- false
	return ctx.SendStatus(fiber.StatusOK)
}

func (h handler) QuitServices(ctx *fiber.Ctx) error {
	h.channels.Quit <- false
	h.channels.Quit <- false

	//close(h.channels.Quit)
	return ctx.SendStatus(fiber.StatusOK)
}
