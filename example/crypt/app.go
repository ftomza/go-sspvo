package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/ftomza/go-sspvo/client"
	"github.com/ftomza/go-sspvo/crypto"
	"github.com/ftomza/go-sspvo/message"
	"github.com/go-resty/resty/v2"
)

var (
	key = `
-----BEGIN PRIVATE KEY-----
MEgCAQAwHwYIKoUDBwEBBgEwEwYHKoUDAgIkAAYIKoUDBwEBAgIEIgQgAnLfE4VXwFTuD5HbBX84W9f/NLDcxNXUWHB+Atu/
6BE=
-----END PRIVATE KEY-----
`
	cert = `
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
`
)

func main() {
	restyClient := resty.New()
	restyClient.SetHostURL("http://85.142.162.12:8031")
	sspvoClient, err := client.NewRestyClient(restyClient,
		client.SetAPIBase("/api"),
		client.SetOGRN("test"),
		client.SetKPP("test"),
	)
	if err != nil {
		log.Fatal(err)
	}

	gostCrypto, _ := crypto.NewGostCrypto(crypto.SetCert(cert), crypto.SetKey(key))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	data, err := sspvoClient.Send(ctx, message.NewActionMessage(
		gostCrypto,
		message.ActionAdd,
		message.DatatypeSubdivisionOrg,
		[]byte(`
<?xml version="1.0" encoding="utf-8"?>
<PackageData>
	<SubdivisionOrg>
		<UID>TEST69</UID>
		<Name>Подвал</Name>
	</SubdivisionOrg>
</PackageData>
	`))).Data()
	cancel()
	if err != nil {
		log.Fatal(err, string(data))
	}

	log.Println("Add SubdivisionOrg:", string(data)) // 2020/10/01 21:09:06 Add SubdivisionOrg: {"IDJWT":"1405161"}

	newMessage := struct {
		IDJWT string `json:"IDJWT"`
	}{}

	err = json.Unmarshal(data, &newMessage)
	if err != nil {
		log.Fatal(err, string(data))
	}

	idJWT, _ := strconv.Atoi(newMessage.IDJWT)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Hour)
	data, err = sspvoClient.Send(ctx, message.NewInfoMessage(gostCrypto, idJWT)).Data()
	cancel()
	if err != nil {
		log.Fatal(err, string(data))
	}

	log.Println("Info:", string(data)) //2020/10/01 21:09:07 Info: {"ResponseToken":".."}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Hour)
	data, err = sspvoClient.Send(ctx, message.NewConfirmMessage(gostCrypto, idJWT)).Data()
	cancel()
	if err != nil {
		log.Fatal(err, string(data))
	}

	log.Println("Confirm:", string(data)) //2020/10/01 21:09:08 Confirm: {"IDJWT":"1405161","Result":"true"}
}
