/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/ftomza/go-sspvo"
)

const (
	pathMethodAction  = "token/new"
	pathMethodCLS     = "cls/request"
	pathMethodInfo    = "token/info"
	pathMethodConfirm = "token/confirm"
)

func SetOGRN(ogrn string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldOGRN] = ogrn
	}
}

func SetKPP(kpp string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldKPP] = kpp
	}
}

func setCert(cert string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		block, _ := pem.Decode([]byte(cert))
		if block != nil {
			f[sspvo.FieldCert] = base64.StdEncoding.EncodeToString(block.Bytes)
			return
		}
		f[sspvo.FieldCert] = cert
	}
}

func setToken(token *sspvo.Token) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldToken] = fmt.Sprintf("%s.%s.%s", token.Header, token.Payload, token.Sign)
	}
}

func setIdJWT(idJWT int) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldIdJWT] = idJWT
	}
}

func setCLS(cls CLS) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldCLS] = cls.String()
	}
}

func setAction(action Action) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldAction] = action.String()
	}
}

func setDatatype(dataType Datatype) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldDataType] = dataType.String()
	}
}

type Message struct {
	Fields sspvo.JWTFields
}

func (m *Message) UpdateJWTFields(fields ...sspvo.Fields) sspvo.Message {
	for _, field := range fields {
		field(m.Fields)
	}
	return m
}

func (m *Message) Init() {
	m.Fields = sspvo.JWTFields{}
}

func (m *Message) PathMethod() string {
	panic("not implement")
}

func (m *Message) GetJWT() ([]byte, error) {
	return json.Marshal(m.Fields)
}

func (m *Message) Response() sspvo.Response {

	panic("not implement")
}

type SignMessage struct {
	Message
	crypto sspvo.Crypto
	data   []byte
}

func (m *SignMessage) UpdateJWTFields(fields ...sspvo.Fields) sspvo.Message {
	m.Message.UpdateJWTFields(fields...)
	return m
}

func (m *SignMessage) GetJWT() ([]byte, error) {
	token, err := m.signToken()
	if err != nil {
		return nil, fmt.Errorf("SignMessage: %w", err)
	}
	fields := sspvo.JWTFields{}
	setToken(token)(fields)
	return json.Marshal(fields)
}

func (m *SignMessage) Init(crypto sspvo.Crypto, data []byte) {
	m.Message.Init()
	m.crypto = crypto
	m.data = data
}

func (m *SignMessage) signToken() (*sspvo.Token, error) {
	m.UpdateJWTFields(setCert(m.crypto.GetCert()))

	jsonHeader, err := json.Marshal(m.Fields)
	if err != nil {
		return nil, err
	}
	header := base64.StdEncoding.EncodeToString(jsonHeader)

	payload := ""
	if m.data != nil {
		payload = base64.StdEncoding.EncodeToString(m.data)
	}

	data4sign := []byte(fmt.Sprintf("%s.%s", header, payload))
	digest := m.crypto.Hash(data4sign)
	signDigest, err := m.crypto.Sign(digest)
	if err != nil {
		return nil, err
	}
	if ok, err := m.crypto.Verify(signDigest, digest); err != nil {
		return nil, err
	} else if !ok {
		return nil, sspvo.ErrBadSign
	}
	sign := base64.StdEncoding.EncodeToString(signDigest)
	return &sspvo.Token{
		Header:  header,
		Payload: payload,
		Sign:    sign,
	}, nil
}
