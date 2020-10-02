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

type InfoMessage struct {
	SignMessage
}

func NewInfoMessage(crypto sspvo.Crypto, idJWT int) *InfoMessage {
	msg := &InfoMessage{}
	msg.Init(crypto, nil)
	msg.UpdateJWTFields(setIdJWT(idJWT), setAction(actionGetMessage))

	return msg
}

func (m *InfoMessage) PathMethod() string {
	return pathMethodInfo
}

func (m *InfoMessage) Response() sspvo.Response {
	resp := response.NewSignResponse()
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
