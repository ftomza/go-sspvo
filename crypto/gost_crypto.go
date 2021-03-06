/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package crypto

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"hash"

	gost_crypto "github.com/ftomza/go-gost-crypto"
	"github.com/ftomza/go-sspvo"

	"github.com/ftomza/gogost/gost3410"
)

//GostCrypto Crypto structure that implements the sspvo.Crypto interface using the crypto package gost3410.
type GostCrypto struct {
	Crypto
	privateKey *gost3410.PrivateKey
	publicKey  *gost3410.PublicKey
}

//NewGostCrypto Creating a new GostCrypto, supports the following options: SetCert, SetKey(Optional).
func NewGostCrypto(opts ...Option) (sspvo.Crypto, error) {
	crypto, err := NewCrypto(opts...)
	if err != nil {
		return nil, err
	}

	publicKey, err := parsePublicKey(crypto.opts.cert)
	if err != nil {
		return nil, fmt.Errorf("gost_crypto: %w", err)
	}

	crypto.hash, err = parseHash(crypto.opts.cert)
	if err != nil {
		return nil, fmt.Errorf("gost_crypto: %w", err)
	}

	var privateKey *gost3410.PrivateKey
	if crypto.opts.key != "" {
		privateKey, err = parsePrivateKey(crypto.opts.key, publicKey)
		if err != nil {
			return nil, fmt.Errorf("gost_crypto: %w", err)
		}
	}

	return &GostCrypto{
		Crypto:     crypto,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func parsePrivateKey(key string, pub *gost3410.PublicKey) (*gost3410.PrivateKey, error) {
	der, err := gost_crypto.DerDecode([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("privateKey: %w", err)
	}

	privateKey, err := gost_crypto.ParsePKCS8PrivateKey(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("privateKey: %w", err)
	}

	publicKey, err := privateKey.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("privateKey: extract PublicKey: %w", err)
	}

	if !bytes.Equal(publicKey.Raw(), pub.Raw()) {
		return nil, errors.New("privateKey: PublicKey not equal")
	}

	return privateKey, nil
}

func parsePublicKeyPartOnCert(cert string) (res []byte, err error) {
	der, err := gost_crypto.DerDecode([]byte(cert))
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	res, err = gost_crypto.ParseCertificate(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	return
}

func parsePublicKey(cert string) (*gost3410.PublicKey, error) {
	pub, err := parsePublicKeyPartOnCert(cert)
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	publicKey, err := gost_crypto.ParsePKIXPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	return publicKey, nil
}

func parseHash(cert string) (hash.Hash, error) {
	pub, err := parsePublicKeyPartOnCert(cert)
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	hashFunc, err := gost_crypto.ParsePKIXPublicKeyHash(pub)
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	return hashFunc, nil
}

//GetVerifyCrypto Get crypto instance for verify by cert
func (c *GostCrypto) GetVerifyCrypto(cert string) (sspvo.Crypto, error) {
	return NewGostCrypto(SetCert(cert))
}

//Hash Get a hash(digest) based on the subtleties of GOST
func (c GostCrypto) Hash(data []byte) (hash []byte) {
	hash = c.Crypto.Hash(data)
	gost_crypto.Reverse(hash)
	return
}

//Sign the digest using the private key generated during initialization
func (c GostCrypto) Sign(digest []byte) (sign []byte, err error) {
	if c.privateKey == nil {
		return nil, fmt.Errorf("gost_crypto: Private Key not set")
	}
	sign, err = c.privateKey.SignDigest(digest, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("gost_crypto: %w", err)
	}
	return
}

//Verify the digest signature with the public key generated at the time of initialization
func (c GostCrypto) Verify(sign, digest []byte) (ok bool, err error) {
	ok, err = c.publicKey.VerifyDigest(digest, sign)
	if err != nil {
		return false, fmt.Errorf("gost_crypto: %w", err)
	}
	return
}
