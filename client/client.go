package client

import (
	"errors"
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/message"
)

type options struct {
	apiBase string
	ogrn    string
	kpp     string
}

type Option func(*options)

func SetAPIBase(apiBase string) Option {
	return func(o *options) {
		o.apiBase = apiBase
	}
}

func SetOGRN(ogrn string) Option {
	return func(o *options) {
		o.ogrn = ogrn
	}
}

func SetKPP(kpp string) Option {
	return func(o *options) {
		o.kpp = kpp
	}
}

type Client struct {
	opts *options
}

func NewClient(opts ...Option) (Client, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.ogrn == "" {
		return Client{}, errors.New("client: ogrn not set")
	}

	if o.kpp == "" {
		return Client{}, errors.New("client: kpp not set")
	}

	return Client{
		opts: &o,
	}, nil
}

func (c *Client) prepareBody(msg sspvo.Message) ([]byte, error) {

	body, err := msg.UpdateJWTFields(message.SetKPP(c.opts.kpp), message.SetOGRN(c.opts.ogrn)).GetJWT()
	if err != nil {
		return nil, err
	}
	return body, err
}
