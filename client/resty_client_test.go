/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ftomza/go-sspvo"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	client *resty.Client
	server *httptest.Server
}

func (suite *ClientTestSuite) SetupTest() {
	suite.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	}))
	suite.client = resty.New()
	suite.client.SetHostURL(suite.server.URL)
}

func Test_ClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TearDownTest() {
	suite.server.Close()
}

func (suite *ClientTestSuite) TestNewRestyClient() {
	type args struct {
		rest *resty.Client
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    sspvo.Client
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				rest: suite.client,
				opts: []Option{SetKPP("KPP"), SetOGRN("OGRN"), SetAPIBase("BASE")},
			},
			want: &RestyClient{
				Client: Client{
					opts: &options{
						apiBase: "BASE",
						ogrn:    "OGRN",
						kpp:     "KPP",
					},
				},
				rest: suite.client,
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				rest: nil,
				opts: []Option{SetOGRN("OGRN"), SetAPIBase("BASE")},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			got, err := NewRestyClient(tt.args.rest, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRestyClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRestyClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *ClientTestSuite) TestRestyClient_Send() {

	badCtx, cancel := context.WithCancel(context.Background())
	cancel()
	badClient := resty.New()
	badClient.SetHostURL("http://bad")
	type fields struct {
		Client Client
		rest   *resty.Client
	}
	type args struct {
		ctx context.Context
		msg sspvo.Message
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantRes    sspvo.Response
		wantResErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Client: Client{
					opts: &options{},
				},
				rest: suite.client,
			},
			args: args{
				ctx: context.Background(),
				msg: Message{false},
			},
			wantRes:    &Response{false, nil},
			wantResErr: false,
		},
		{
			name: "fail prepare",
			fields: fields{
				Client: Client{
					opts: &options{},
				},
				rest: suite.client,
			},
			args: args{
				ctx: context.Background(),
				msg: Message{true},
			},
			wantRes:    &Response{true, nil},
			wantResErr: true,
		},
		{
			name: "fail resty",
			fields: fields{
				Client: Client{
					opts: &options{},
				},
				rest: badClient,
			},
			args: args{
				ctx: badCtx,
				msg: Message{false},
			},
			wantRes:    &Response{false, nil},
			wantResErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			c := &RestyClient{
				Client: tt.fields.Client,
				rest:   tt.fields.rest,
			}
			gotRes := c.Send(tt.args.ctx, tt.args.msg)
			if (gotRes.Error() != nil) != tt.wantResErr {
				t.Errorf("Send().Error() = %v, wantErr %v", gotRes.Error(), tt.wantResErr)
				return
			}
			gotRes.SetError(nil)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Send() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func init() {

}
