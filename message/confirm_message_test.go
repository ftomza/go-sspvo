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

func TestNewConfirmMessage(t *testing.T) {
	type args struct {
		crypto sspvo.Crypto
		idJWT  int
	}
	tests := []struct {
		name string
		args args
		want *ConfirmMessage
	}{
		{
			name: "ok",
			args: args{
				idJWT: 1,
			},
			want: &ConfirmMessage{
				SignMessage: SignMessage{
					Message: Message{
						Fields: sspvo.JWTFields{sspvo.FieldAction: actionMessageConfirm.String(), sspvo.FieldIdJWT: 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfirmMessage(tt.args.crypto, tt.args.idJWT); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfirmMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfirmMessage_PathMethod(t *testing.T) {
	type fields struct {
		SignMessage SignMessage
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			want: "token/confirm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ConfirmMessage{
				SignMessage: tt.fields.SignMessage,
			}
			if got := m.PathMethod(); got != tt.want {
				t.Errorf("PathMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
