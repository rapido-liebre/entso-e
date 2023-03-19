package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ParseQueryParams(ctx *fiber.Ctx, expectedParamsCount int) ([]string, error) {
	params := strings.Split(string(ctx.Request().URI().QueryString()), "=")
	if len(params) != expectedParamsCount {
		return []string{}, fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("Invalid request params %s", ctx.Request().URI().QueryString()))
	}
	return params, nil
}

/*
func (h handler) AuthenticateUser(ctx *fiber.Ctx) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//TODO parse username credentials in query params
	//if params, err := ParseQueryParams(ctx, 2); err == nil {
	//	fmt.Println(cfg.Params.KerberosDir, params)
	//}
	fmt.Println(cfg.Params.KerberosFile)

	// Authenticate user
	keytab, err := keytab2.Load(cfg.Params.KerberosFile) //(*keytab.Keytab, error)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data, err := keytab.Marshal()
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Println(data)
	//for entry := range keytab.Entries {
	//	fmt.Println(entry)
	//}
	fmt.Println(keytab.Entries)

	//encryptionKey, err := keytab.GetEncryptionKey("HTTP/wams-dev.inno.pse.pl@AD.PSE.PL", 3) //princName types.PrincipalName, realm string, kvno int, etype int32) // (types.EncryptionKey, error)
	//fmt.Println(encryptionKey)

	//if config.UseMinio() {
	//	mc, _ := config.GetMinioClient()
	//	buckets, err = mc.ListBuckets()
	//} else {
	//	if !config.UseMinioTestData() {
	//		return fiber.NewError(fiber.StatusInternalServerError, "No minIO client defined")
	//	}
	//	buckets, err = config.GetFakeBuckets()
	//}
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	//for _, bucket := range buckets {
	//	fmt.Println(bucket)
	//}

	return ctx.SendStatus(fiber.StatusOK)
}
*/
