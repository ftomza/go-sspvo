/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package crypto

import (
	"errors"
	"hash"
)

type options struct {
	cert string
	key  string
}

type Option func(*options)

func SetCert(cert string) Option {
	return func(o *options) {
		o.cert = cert
	}
}

func SetKey(key string) Option {
	return func(o *options) {
		o.key = key
	}
}

type Crypto struct {
	opts *options
	hash hash.Hash
}

func NewCrypto(opts ...Option) (Crypto, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.cert == "" {
		return Crypto{}, errors.New("crypto: cert not set")
	}

	return Crypto{
		opts: &o,
	}, nil
}

func (c *Crypto) GetCert() string {
	return c.opts.cert
}

func (c *Crypto) Hash(data []byte) (hash []byte) {
	hashes := c.hash
	_, _ = hashes.Write(data)
	hash = hashes.Sum(nil)
	return
}
