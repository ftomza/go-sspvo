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

type ConfirmMessage struct {
	SignMessage
}

func NewConfirmMessage(crypto sspvo.Crypto, idJWT int) *ConfirmMessage {
	msg := &ConfirmMessage{}
	msg.Init(crypto, nil)
	msg.UpdateJWTFields(setIdJWT(idJWT), setAction(actionMessageConfirm))

	return msg
}

func (m *ConfirmMessage) PathMethod() string {
	return pathMethodConfirm
}

func (m *ConfirmMessage) Response() sspvo.Response {
	resp := response.NewSignResponse()
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
