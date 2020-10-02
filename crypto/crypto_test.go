/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package crypto

import (
	"crypto/md5"
	"hash"
	"reflect"
	"testing"
)

func TestNewCrypto(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Crypto
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				opts: []Option{SetCert("CERT"), SetKey("KEY")},
			},
			want: Crypto{
				opts: &options{
					cert: "CERT",
					key:  "KEY",
				},
				hash: nil,
			},
			wantErr: false,
		},
		{
			name: "fail cert",
			args: args{
				opts: []Option{SetKey("KEY")},
			},
			want: Crypto{
				opts: nil,
				hash: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCrypto(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrypto() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypto_Hash(t *testing.T) {
	hasher := md5.New()
	type fields struct {
		opts *options
		hash hash.Hash
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHash []byte
	}{
		{
			name: "ok",
			fields: fields{
				hash: hasher,
			},
			args: args{
				data: []byte("test"),
			},
			wantHash: []byte{9, 143, 107, 205, 70, 33, 211, 115, 202, 222, 78, 131, 38, 39, 180, 246},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crypto{
				opts: tt.fields.opts,
				hash: tt.fields.hash,
			}
			if gotHash := c.Hash(tt.args.data); !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("Hash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
