package client

import (
	"context"
	"fmt"
	"github.com/ftomza/go-sspvo"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	Client

	rest *resty.Client
}

func NewRestyClient(rest *resty.Client, opts ...Option) (sspvo.Client, error) {
	client, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	rest.SetHeader("Content-Type", "application/json")
	return &RestyClient{
		Client: client,
		rest:   rest,
	}, nil
}

func (c *RestyClient) Send(ctx context.Context, msg sspvo.Message) (res sspvo.Response) {

	res = msg.Response()

	body, err := c.prepareBody(msg)
	if err != nil {
		res.SetError(fmt.Errorf("client prepare body: %w", err))
		return
	}

	resp, err := c.rest.R().
		SetContext(ctx).
		SetBody(body).
		Post(fmt.Sprintf("%s/%s", c.opts.apiBase, msg.PathMethod()))

	if err != nil {
		res.SetError(fmt.Errorf("client send: %w", err))
		return
	}
	res.SetClientResponse(&sspvo.ClientResponse{
		Code:   resp.StatusCode(),
		Body:   resp.Body(),
		Header: resp.Header(),
	})
	return
}
