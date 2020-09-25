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

type GostCrypto struct {
	Crypto
	privateKey *gost3410.PrivateKey
	publicKey  *gost3410.PublicKey
}

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

func parsePublicKey(cert string) (*gost3410.PublicKey, error) {
	der, err := gost_crypto.DerDecode([]byte(cert))
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	publicKey, err := gost_crypto.ParsePKIXPublicKey(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("publicKey: %w", err)
	}

	return publicKey, nil
}

func parseHash(cert string) (hash.Hash, error) {
	der, err := gost_crypto.DerDecode([]byte(cert))
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	hashFunc, err := gost_crypto.ParsePKIXPublicKeyHash(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	return hashFunc, nil
}

func (c *GostCrypto) GetVerifyCrypto(cert string) (sspvo.Crypto, error) {
	return NewGostCrypto(SetCert(cert))
}

func (c *GostCrypto) Hash(data []byte) (hash []byte) {
	hash = c.Crypto.Hash(data)
	gost_crypto.Reverse(hash)
	return
}

func (c *GostCrypto) Sign(digest []byte) (sign []byte, err error) {
	if c.privateKey == nil {
		return nil, fmt.Errorf("gost_crypto: Private Key not set")
	}
	sign, err = c.privateKey.SignDigest(digest, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("gost_crypto: %w", err)
	}
	return
}

func (c *GostCrypto) Verify(sign, digest []byte) (ok bool, err error) {
	ok, err = c.publicKey.VerifyDigest(digest, sign)
	if err != nil {
		return false, fmt.Errorf("gost_crypto: %w", err)
	}
	return
}
