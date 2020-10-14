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

//SetCert To set the cert option. Assigns a certificate in the DER x509 format.
func SetCert(cert string) Option {
	return func(o *options) {
		o.cert = cert
	}
}

//SetKey To set the key option. Assigns a Private Key in the DER PKCS8 format.
func SetKey(key string) Option {
	return func(o *options) {
		o.key = key
	}
}

//Crypto Basic structure that implements the sspvo.Crypto interface.
type Crypto struct {
	opts *options
	hash hash.Hash
}

//NewCrypto Creating a new base Crypto, supports the following options: SetCert, SetKey(Optional).
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

//GetCert Receive options cert
func (c *Crypto) GetCert() string {
	return c.opts.cert
}

//Hash Get hash(digest) of data
func (c *Crypto) Hash(data []byte) (hash []byte) {
	hashes := c.hash
	hashes.Reset()
	_, _ = hashes.Write(data)
	hash = hashes.Sum(nil)
	return
}
