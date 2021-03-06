/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package client

import (
	"context"
	"fmt"

	"github.com/ftomza/go-sspvo"

	"github.com/go-resty/resty/v2"
)

//RestyClient Client structure that implements the sspvo.Client interface using the rest client resty.Client.
type RestyClient struct {
	Client

	rest *resty.Client
}

//NewRestyClient Creating a new RestyClient, supports the following options: SetOGRN, SetKPP, SetAPIBase and instance resty.Client.
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

//Send a instance message, see message.Message, with the specified context and return an instance of the prepared response based on the interface sspvo.Response
func (c *RestyClient) Send(ctx context.Context, msg sspvo.Message) (res sspvo.Response) {

	res = msg.Response()

	body, err := c.PrepareBody(msg)
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
