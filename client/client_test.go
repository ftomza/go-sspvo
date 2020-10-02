/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package client

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ftomza/go-sspvo"
)

type Message struct {
	fail bool
}

func (m Message) UpdateJWTFields(fields ...sspvo.Fields) sspvo.Message {
	return m
}

func (m Message) PathMethod() string {
	return ""
}

func (m Message) GetJWT() ([]byte, error) {
	if m.fail {
		return nil, errors.New("fail")
	} else {
		return []byte("ok"), nil
	}
}

func (m Message) Response() sspvo.Response {
	return &Response{m.fail, nil}
}

type Response struct {
	fail bool
	err  error
}

func (r Response) ClientResponse() *sspvo.ClientResponse {
	panic("implement me")
}

func (r Response) SetClientResponse(resp *sspvo.ClientResponse) {
}

func (r *Response) SetError(err error) {
	r.err = err
}

func (r Response) Error() error {
	return r.err
}

func (r Response) Data() ([]byte, error) {
	panic("implement me")
}

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				opts: []Option{SetKPP("KPP"), SetOGRN("OGRN"), SetAPIBase("BASE")},
			},
			want: Client{
				opts: &options{
					apiBase: "BASE",
					ogrn:    "OGRN",
					kpp:     "KPP",
				},
			},
			wantErr: false,
		},
		{
			name: "fail kpp",
			args: args{
				opts: []Option{SetOGRN("OGRN"), SetAPIBase("BASE")},
			},
			want:    Client{},
			wantErr: true,
		},
		{
			name: "fail ogrn",
			args: args{
				opts: []Option{SetKPP("KPP")},
			},
			want:    Client{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_prepareBody(t *testing.T) {
	type fields struct {
		opts *options
	}
	type args struct {
		msg sspvo.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				opts: &options{
					apiBase: "",
					ogrn:    "",
					kpp:     "",
				},
			},
			args: args{
				msg: Message{false},
			},
			want:    []byte("ok"),
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				opts: &options{
					apiBase: "",
					ogrn:    "",
					kpp:     "",
				},
			},
			args: args{
				msg: Message{true},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				opts: tt.fields.opts,
			}
			got, err := c.prepareBody(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
