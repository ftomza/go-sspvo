/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package test_server_epgu

import (
	"fmt"
	"log"

	"github.com/ftomza/go-sspvo"

	"github.com/ftomza/go-sspvo/crypto"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Messages     map[int]interface{}
	Port         string
	OGRN         string
	KPP          string
	CertOGRN     string
	Cert         string
	Key          string
	crypto       sspvo.Crypto
	cryptoVerify sspvo.Crypto
}

func (s *Server) apiClsRequest(ctx *fiber.Ctx) error {
	clsBody := struct {
		OGRN string `json:"OGRN"`
		KPP  string `json:"KPP"`
		CLS  string `json:"CLS"`
	}{}
	err := ctx.BodyParser(&clsBody)
	if err != nil {
		return err
	}
	err = s.checkAuth(clsBody.OGRN, clsBody.KPP)
	if err != nil {
		return err
	}

	if clsBody.CLS == "LevelBudget" {
		return ctx.Send([]byte(`
<?xml version="1.0" encoding="UTF-8"?>
<LevelBudget>
<Budget><ID>1</ID><Code></Code><Name>Федеральный</Name><Actual>true</Actual></Budget>
<Budget><ID>2</ID><Code></Code><Name>Региональный</Name><Actual>true</Actual></Budget>
<Budget><ID>3</ID><Code></Code><Name>Муниципальный</Name><Actual>true</Actual></Budget>
</LevelBudget>
`))
	}

	return fiber.NewError(fiber.StatusNotFound,
		`{"Error":"{\"0\":\"field: OriginalErrorText message :Неизвестный тип классификатора\"}"}`)
}

func (s *Server) apiTokenInfo(ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
}

func (s *Server) apiTokenConfirm(ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
}

func (s *Server) apiTokenNew(ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
}

func (s *Server) checkAuth(ogrn string, kpp string) error {
	if s.OGRN != ogrn || s.KPP != kpp {
		return fiber.NewError(fiber.StatusNotFound,
			fmt.Sprintf("{\"Error\":\"{\\\"0\\\":\\\"field: OriginalErrorText message :Организация с OGRN \\\\\\\"%s\\\\\\\" и КПП \\\\\\\"%s\\\\\\\" не найдена\\\"}\"}", ogrn, kpp))
	}
	return nil
}

func RunServerDefault() {
	RunServer(&Server{
		Messages: map[int]interface{}{},
		Port:     "7777",
		OGRN:     "test",
		KPP:      "test",
		CertOGRN: `
-----BEGIN CERTIFICATE-----
MIIEfDCCBCmgAwIBAgIEXek0LjAKBggqhQMHAQEDAjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wHhcNMjAwOTIyMjEw
MDAwWhcNNDAwOTIyMjEwMDAwWjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Qn9C10YLQtdGA0LHR
g9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0GA1UECwwW0KDRg9C6
0L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7MDkGA1UEAwwy0JDQ
u9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wZjAfBggqhQMHAQEBATATBgcqhQMCAiQA
BggqhQMHAQECAgNDAARAyuHXvOdPT/R94KICw82bdgiBfEXkEJxqXIN4uav8zIvgDe/q7yzK+HJnbLWLIWc2z+eqbaiUbj0Y
e1RoNUa5NaOCAZ4wggGaMA4GA1UdDwEB/wQEAwIB/jAxBgNVHSUEKjAoBggrBgEFBQcDAQYIKwYBBQUHAwIGCCsGAQUFBwMD
BggrBgEFBQcDBDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSalTlfa+t/MpLv76stCkVlU18TazCCASMGA1UdIwSCARow
ggEWgBSalTlfa+t/MpLv76stCkVlU18Ta6GB96SB9DCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L2CBF3pNC4wCgYIKoUD
BwEBAwIDQQBlY4HdS/G7zAWOEWH6pBx4FSli5ipbEtvr/lkjEApvlrch5cMlmy7rglAbE7ct+sKFtDKv6cIhqu3rQMAla/gb
-----END CERTIFICATE-----
`,
		Key: `
-----BEGIN PRIVATE KEY-----
MEgCAQAwHwYIKoUDBwEBBgEwEwYHKoUDAgIkAAYIKoUDBwEBAgIEIgQgJuknaNfag4otfbYpIFTMF815j5SSXnvd3NyWdXPk
E5Y=
-----END PRIVATE KEY-----
`,
		Cert: `
-----BEGIN CERTIFICATE-----
MIIEfDCCBCmgAwIBAgIEGh2orzAKBggqhQMHAQEDAjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wHhcNMjAxMDE1MjEw
MDAwWhcNNDAxMDE1MjEwMDAwWjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Qn9C10YLQtdGA0LHR
g9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0GA1UECwwW0KDRg9C6
0L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7MDkGA1UEAwwy0JDQ
u9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wZjAfBggqhQMHAQEBATATBgcqhQMCAiQA
BggqhQMHAQECAgNDAARAq/L5/D3TSZjQbBS/mkSGgMsudv33fycA0VOwF12HFp7JRZW2TJes5YZbevgA2AhqJT48pT6niIYC
ezJMT4a+MaOCAZ4wggGaMA4GA1UdDwEB/wQEAwIB/jAxBgNVHSUEKjAoBggrBgEFBQcDAQYIKwYBBQUHAwIGCCsGAQUFBwMD
BggrBgEFBQcDBDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRFWyEw+Ry6/S+kgwnhnT7fGtZqUDCCASMGA1UdIwSCARow
ggEWgBRFWyEw+Ry6/S+kgwnhnT7fGtZqUKGB96SB9DCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L2CBBodqK8wCgYIKoUD
BwEBAwIDQQDgarmfivpl0uMY60AJDfyDEu41P5jqrj6jMTt+pBOZ036ROjb8ao7UpAtrP8mMi9MH68XZirxmab6NEshDjTC7
-----END CERTIFICATE-----
`,
	})
}

func RunServer(server *Server) {

	var (
		err error
	)

	server.crypto, err = crypto.NewGostCrypto(crypto.SetCert(server.Cert), crypto.SetKey(server.Key))
	if err != nil {
		log.Fatal(err)
	}

	server.cryptoVerify, err = crypto.NewGostCrypto(crypto.SetCert(server.CertOGRN))
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	api := app.Group("/api")
	api.Post("/cls/request", server.apiClsRequest)

	token := api.Group("/token")
	token.Post("/info", server.apiTokenInfo)
	token.Post("/confirm", server.apiTokenConfirm)
	token.Post("/new", server.apiTokenNew)

	err = app.Listen(fmt.Sprintf(":%s", server.Port))
	if err != nil {
		log.Fatal(err)
	}
}
