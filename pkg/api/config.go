package api

import (
	"entso-e_reports/pkg/common/config"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h handler) GetConfig(ctx *fiber.Ctx) error {
	// load config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Println("Failed at loading config", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	for _, path := range config.GetConfigPaths() {
		cfgPath := strings.Join([]string{path, config.GetConfigFilename()}, "/")
		if _, err := os.Stat(cfgPath); err == nil {
			cfg.Path = cfgPath
		}
	}

	//ss := "{Params:\n   {\n    TimeInterval: 5,\n    WarningSize: 10,\n    RedAlertSize: 3,\n    InputDir: \"/Users/rapido_liebre/GolandProjects/wams_archiver/tests/input\",\n    OutputDir: \"/Users/rapido_liebre/GolandProjects/wams_archiver/tests/output\",\n    Port: \":3055\"\n    },\n    Path:\"./\"\n}"
	return ctx.Status(fiber.StatusOK).JSON(cfg)
}

func (h handler) UpdateConfig(ctx *fiber.Ctx) error {
	// load config
	cfg, err := config.GetConfig()
	if err != nil {
		log.Println("Failed at loading config", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// validate all fields
	if err = cfg.Validate(); err != nil {
		log.Println("Failed at validating config", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// send config to main loop
	h.channels.CfgUpdate <- cfg

	return ctx.Status(fiber.StatusOK).JSON(cfg)
}
