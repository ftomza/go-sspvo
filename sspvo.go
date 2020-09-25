/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package sspvo

import (
	"context"
	"errors"
	"net/http"
)

const (
	VERSION = "1.0.0"
)

var (
	ErrBadSign = errors.New("crypto: bad sign")
)

type Crypto interface {
	GetVerifyCrypto(cert string) (Crypto, error)
	GetCert() string
	Hash(data []byte) (hash []byte)
	Sign(digest []byte) (sign []byte, err error)
	Verify(sign, digest []byte) (ok bool, err error)
}

type Client interface {
	Send(ctx context.Context, msg Message) (res Response)
}

type Token struct {
	Header  string
	Payload string
	Sign    string
}

type FieldName string

const (
	FieldOGRN     FieldName = "OGRN"
	FieldKPP      FieldName = "KPP"
	FieldCLS      FieldName = "CLS"
	FieldCert     FieldName = "Cert64"
	FieldToken    FieldName = "token"
	FieldAction   FieldName = "action"
	FieldDataType FieldName = "data_type"
	FieldIdJWT    FieldName = "IDJWT"
)

var AllField = []FieldName{
	FieldOGRN,
	FieldKPP,
	FieldCLS,
	FieldCert,
	FieldToken,
	FieldAction,
	FieldDataType,
	FieldIdJWT,
}

func (e FieldName) IsValid() bool {
	switch e {
	case FieldOGRN,
		FieldKPP,
		FieldCLS,
		FieldCert,
		FieldToken,
		FieldAction,
		FieldDataType,
		FieldIdJWT:
		return true
	}
	return false
}

func (e FieldName) String() string {
	return string(e)
}

type JWTFields map[FieldName]string

type Fields func(JWTFields)

type Message interface {
	UpdateJWTFields(fields ...Fields) Message
	PathMethod() string
	GetJWT() ([]byte, error)
	Response() Response
}

type ClientResponse struct {
	Code   int
	Body   []byte
	Header http.Header
}

type Response interface {
	ClientResponse() *ClientResponse
	SetClientResponse(resp *ClientResponse)
	SetError(err error)
	Error() error
	Data() ([]byte, error)
}
