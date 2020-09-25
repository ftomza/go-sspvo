/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type ActionMessage struct {
	SignMessage
}

func NewActionMessage(crypto sspvo.Crypto, action, dataType string, data []byte) *ActionMessage {
	msg := &ActionMessage{}
	msg.Init(crypto, data)
	msg.UpdateJWTFields(setAction(action), setDataType(dataType))

	return msg
}

func (m *ActionMessage) PathMethod() string {
	return pathMethodAction
}

func (m *ActionMessage) Response() sspvo.Response {
	resp := response.NewActionResponse(m.Fields[sspvo.FieldAction], m.Fields[sspvo.FieldDataType])
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
