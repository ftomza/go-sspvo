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

func TestAction_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Action
		want bool
	}{
		{
			name: "ok",
			e:    "Add",
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

func TestAction_String(t *testing.T) {
	tests := []struct {
		name string
		e    Action
		want string
	}{
		{
			name: "ok",
			e:    "Add",
			want: "Add",
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

func TestDatatype_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Datatype
		want bool
	}{
		{
			name: "ok",
			e:    "ege",
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

func TestDatatype_String(t *testing.T) {
	tests := []struct {
		name string
		e    Datatype
		want string
	}{
		{
			name: "ok",
			e:    "ege",
			want: "ege",
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

func TestNewActionMessage(t *testing.T) {
	type args struct {
		crypto   sspvo.Crypto
		action   Action
		datatype Datatype
		data     []byte
	}
	tests := []struct {
		name string
		args args
		want *ActionMessage
	}{
		{
			name: "ok",
			args: args{
				action:   ActionAdd,
				datatype: DatatypeSubdivisionOrg,
			},
			want: &ActionMessage{
				SignMessage: SignMessage{
					Message: Message{
						Fields: sspvo.JWTFields{sspvo.FieldAction: ActionAdd.String(), sspvo.FieldDataType: DatatypeSubdivisionOrg.String()},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewActionMessage(tt.args.crypto, tt.args.action, tt.args.datatype, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewActionMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionMessage_PathMethod(t *testing.T) {
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
			want: "token/new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ActionMessage{
				SignMessage: tt.fields.SignMessage,
			}
			if got := m.PathMethod(); got != tt.want {
				t.Errorf("PathMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
