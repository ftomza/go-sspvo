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

type CLSMessage struct {
	Message
}

func NewCLSMessage(cls string) *CLSMessage {
	msg := &CLSMessage{}
	msg.Init()
	msg.UpdateJWTFields(setCLS(cls))

	return msg
}

func (m *CLSMessage) PathMethod() string {
	return pathMethodCLS
}

func (m *CLSMessage) Response() sspvo.Response {
	return response.NewCLSResponse(m.Fields[sspvo.FieldCLS])
}
