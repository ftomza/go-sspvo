/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

import (
	"reflect"
	"testing"

	"github.com/ftomza/go-sspvo"
)

func TestInfoAllMessage_PathMethod(t *testing.T) {
	type fields struct {
		Message Message
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			want: "token/info",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InfoAllMessage{
				Message: tt.fields.Message,
			}
			if got := m.PathMethod(); got != tt.want {
				t.Errorf("PathMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInfoAllMessage(t *testing.T) {
	tests := []struct {
		name string
		want *InfoAllMessage
	}{
		{
			name: "ok",
			want: &InfoAllMessage{
				Message: Message{
					Fields: sspvo.JWTFields{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInfoAllMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInfoAllMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
