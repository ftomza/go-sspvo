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

func TestNewCLSMessage(t *testing.T) {
	type args struct {
		cls CLS
	}
	tests := []struct {
		name string
		args args
		want *CLSMessage
	}{
		{
			name: "ok",
			args: args{
				cls: CLSDirections,
			},
			want: &CLSMessage{
				Message: Message{
					Fields: sspvo.JWTFields{sspvo.FieldCLS: CLSDirections.String()},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCLSMessage(tt.args.cls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCLSMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLS_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    CLS
		want bool
	}{
		{
			name: "ok",
			e:    "Subject",
			want: true,
		},
		{
			name: "false",
			e:    "Bad",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLS_String(t *testing.T) {
	tests := []struct {
		name string
		e    CLS
		want string
	}{
		{
			name: "ok",
			e:    "Subject",
			want: "Subject",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLSMessage_PathMethod(t *testing.T) {
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
			want: "cls/request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CLSMessage{
				Message: tt.fields.Message,
			}
			if got := m.PathMethod(); got != tt.want {
				t.Errorf("PathMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
