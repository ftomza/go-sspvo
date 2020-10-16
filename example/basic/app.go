package main

import (
	"context"
	"log"
	"time"

	"github.com/ftomza/go-sspvo/test_server_epgu"

	"github.com/ftomza/go-sspvo/client"
	"github.com/ftomza/go-sspvo/message"

	"github.com/go-resty/resty/v2"
)

func main() {

	go test_server_epgu.RunServerDefault()

	restyClient := resty.New()
	restyClient.SetHostURL("http://localhost:7777")
	restyClient.SetDebug(false)
	sspvoClient, err := client.NewRestyClient(restyClient,
		client.SetAPIBase("/api"),
		client.SetOGRN("test"),
		client.SetKPP("test"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := sspvoClient.Send(ctx, message.NewCLSMessage(message.CLSLevelBudget)).Data()
	if err != nil {
		log.Fatal(err, string(data))
	}

	log.Print(string(data))
}
