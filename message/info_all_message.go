/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

type InfoAllMessage struct {
	Message
}

func NewInfoAllMessage() *InfoAllMessage {
	msg := &InfoAllMessage{}
	msg.Init()

	return msg
}

func (m *InfoAllMessage) PathMethod() string {
	return pathMethodInfo
}
