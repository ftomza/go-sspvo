/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package crypto

import (
	"encoding/hex"
	"hash"
	"math/big"
	"reflect"
	"testing"

	"github.com/ftomza/go-sspvo"

	"github.com/ftomza/gogost/gost34112012256"

	"github.com/ftomza/gogost/gost3410"
)

var (
	validCert = `
-----BEGIN CERTIFICATE-----
MIIEfDCCBCmgAwIBAgIEXek0LjAKBggqhQMHAQEDAjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wHhcNMjAwOTIyMjEw
MDAwWhcNNDAwOTIyMjEwMDAwWjCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Qn9C10YLQtdGA0LHR
g9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0GA1UECwwW0KDRg9C6
0L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7MDkGA1UEAwwy0JDQ
u9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L0wZjAfBggqhQMHAQEBATATBgcqhQMCAiQA
BggqhQMHAQECAgNDAARAyuHXvOdPT/R94KICw82bdgiBfEXkEJxqXIN4uav8zIvgDe/q7yzK+HJnbLWLIWc2z+eqbaiUbj0Y
e1RoNUa5NaOCAZ4wggGaMA4GA1UdDwEB/wQEAwIB/jAxBgNVHSUEKjAoBggrBgEFBQcDAQYIKwYBBQUHAwIGCCsGAQUFBwMD
BggrBgEFBQcDBDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSalTlfa+t/MpLv76stCkVlU18TazCCASMGA1UdIwSCARow
ggEWgBSalTlfa+t/MpLv76stCkVlU18Ta6GB96SB9DCB8TELMAkGA1UEBhMCUlUxKjAoBgNVBAgMIdCh0LDQvdC60YLRii3Q
n9C10YLQtdGA0LHRg9GA0LPRijEuMCwGA1UECgwl0JbRg9GA0L3QsNC7ICLQodC+0LLRgNC10LzQtdC90L3QuNC6IjEfMB0G
A1UECwwW0KDRg9C60L7QstC+0LTRgdGC0LLQvjEoMCYGA1UEDAwf0JPQu9Cw0LLQvdGL0Lkg0YDQtdC00LDQutGC0L7RgDE7
MDkGA1UEAwwy0JDQu9C10LrRgdCw0L3QtNGAINCh0LXRgNCz0LXQtdCy0LjRhyDQn9GD0YjQutC40L2CBF3pNC4wCgYIKoUD
BwEBAwIDQQBlY4HdS/G7zAWOEWH6pBx4FSli5ipbEtvr/lkjEApvlrch5cMlmy7rglAbE7ct+sKFtDKv6cIhqu3rQMAla/gb
-----END CERTIFICATE-----
`
	validKey = `
-----BEGIN PRIVATE KEY-----
MEgCAQAwHwYIKoUDBwEBBgEwEwYHKoUDAgIkAAYIKoUDBwEBAgIEIgQgAnLfE4VXwFTuD5HbBX84W9f/NLDcxNXUWHB+Atu/
6BE=
-----END PRIVATE KEY-----
`
	badCertOrKey = `
-----BEGIN DER-----
QkFE
-----END DER-----
`
	badPub = `
-----BEGIN CERTIFICATE-----
MIIEgzCCBCWgAwIBAgIEUsAnJzAOBgorBgEEAa1ZAQMCBQAwgfExCzAJBgNVBAYTAlJVMSowKAYDVQQIDCHQodCw0L3QutGC
0Yot0J/QtdGC0LXRgNCx0YPRgNCz0YoxLjAsBgNVBAoMJdCW0YPRgNC90LDQuyAi0KHQvtCy0YDQtdC80LXQvdC90LjQuiIx
HzAdBgNVBAsMFtCg0YPQutC+0LLQvtC00YHRgtCy0L4xKDAmBgNVBAwMH9CT0LvQsNCy0L3Ri9C5INGA0LXQtNCw0LrRgtC+
0YAxOzA5BgNVBAMMMtCQ0LvQtdC60YHQsNC90LTRgCDQodC10YDQs9C10LXQstC40Ycg0J/Rg9GI0LrQuNC9MB4XDTIwMDky
OTIxMDAwMFoXDTQwMDkyOTIxMDAwMFowgfExCzAJBgNVBAYTAlJVMSowKAYDVQQIDCHQodCw0L3QutGC0Yot0J/QtdGC0LXR
gNCx0YPRgNCz0YoxLjAsBgNVBAoMJdCW0YPRgNC90LDQuyAi0KHQvtCy0YDQtdC80LXQvdC90LjQuiIxHzAdBgNVBAsMFtCg
0YPQutC+0LLQvtC00YHRgtCy0L4xKDAmBgNVBAwMH9CT0LvQsNCy0L3Ri9C5INGA0LXQtNCw0LrRgtC+0YAxOzA5BgNVBAMM
MtCQ0LvQtdC60YHQsNC90LTRgCDQodC10YDQs9C10LXQstC40Ycg0J/Rg9GI0LrQuNC9MF4wFgYKKwYBBAGtWQEGAgYIKoZI
zj0DAQcDRAAEQQRJ47ST4jDv6emPar8XzCAcIb2adsob+TH53QR7YsJHsX6lFh1Y3zpZnfVc/ehMRbD9UcubR5QMptQcGp6k
7PEto4IBnjCCAZowDgYDVR0PAQH/BAQDAgH+MDEGA1UdJQQqMCgGCCsGAQUFBwMBBggrBgEFBQcDAgYIKwYBBQUHAwMGCCsG
AQUFBwMEMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFMvx8MZSQjkPMV8WNliFVx4we65ZMIIBIwYDVR0jBIIBGjCCARaA
FMvx8MZSQjkPMV8WNliFVx4we65ZoYH3pIH0MIHxMQswCQYDVQQGEwJSVTEqMCgGA1UECAwh0KHQsNC90LrRgtGKLdCf0LXR
gtC10YDQsdGD0YDQs9GKMS4wLAYDVQQKDCXQltGD0YDQvdCw0LsgItCh0L7QstGA0LXQvNC10L3QvdC40LoiMR8wHQYDVQQL
DBbQoNGD0LrQvtCy0L7QtNGB0YLQstC+MSgwJgYDVQQMDB/Qk9C70LDQstC90YvQuSDRgNC10LTQsNC60YLQvtGAMTswOQYD
VQQDDDLQkNC70LXQutGB0LDQvdC00YAg0KHQtdGA0LPQtdC10LLQuNGHINCf0YPRiNC60LjQvYIEUsAnJzAOBgorBgEEAa1Z
AQMCBQADSAAwRQIhAL8w+O7XUmYUQfhaCTF0VLz+mB9NYXXT7TXfVBwMtb5kAiBYZ/XkDSyUHCKUPEOFsIH9XXg2wtN7+Q55
SNoNf4LG8g==
-----END CERTIFICATE-----`
	badKey = `
-----BEGIN PRIVATE KEY-----
MIIBGwIBADCB8AYKKwYBBAGtWQEGAjCB4QIBATAsBgcqhkjOPQEBAiEA/////wAAAAEAAAAAAAAAAAAAAAD/////////////
//8wRQIhAP////8AAAABAAAAAAAAAAAAAAAA///////////////8AiBaxjXYqjqT57PrvVV2mIa8ZR0GsMxTsPY7zjw+J9Jg
SwRBBGsX0fLhLEJH+Lzm5WOkQPJ3A32BLeszoPShOUXYmMKWT+NC4v4af5uO5+tKfA+eFivOM1drMV7Oy7ZAaDe/UfUCIQD/
////AAAAAP//////////vOb6racXnoTzucrC/GMlUQIBAQQjAiEAt4eIOCHC2M7HJaLzKDZJ1XjwbHzmYLiPmVoqdlN82qk=
-----END PRIVATE KEY-----
`
	testHash, _    = hex.DecodeString("57381c88028d0db1d099af299d2b596bcf148707fdf2e5f104551b193808a512")
	signExtern, _  = hex.DecodeString("187c82f8f70620ae217897f49c61b059944faebaebf07f7621272dea77d8af49c86a1135c418e25a4d7612b1f1b7d4ee4b00559a7d7ee6f7c708c41453396b55")
	signExtern2, _ = hex.DecodeString("1068bd702e8f0ff9bfafb61a78f5e7fcbd7b4ded63c6d734daa9c72a13143bd26f7bc9b249e537b04a0b84d7b508a3c6b70b3f50182d361cd050d925997ecd85")
)

func getBigInt(t *testing.T, in, mark string) *big.Int {
	v, ok := new(big.Int).SetString(in, 10)
	if !ok {
		t.Fatal("cannot convert string to bigInt")
	}
	return v
}

func getPublicKey(t *testing.T) *gost3410.PublicKey {
	pub := &gost3410.PublicKey{
		C: gost3410.CurveIdGostR34102001CryptoProXchAParamSet(),
		X: getBigInt(t, "63233666624051439876354823295566418637012564188384438200469674371110357426634", "puk"),
		Y: getBigInt(t, "24299932244005800117978005500793438667981994951685184390218551551204573253088", "puk"),
	}
	println("Pub:", hex.EncodeToString(pub.Raw()))
	return pub
}
func getPublicKey2(t *testing.T) *gost3410.PublicKey {
	return &gost3410.PublicKey{
		C: gost3410.CurveIdGostR34102001CryptoProXchAParamSet(),
		X: getBigInt(t, "63233666624051439876354823295566418637012564188384438200469674371110357426634", "puk2"),
		Y: getBigInt(t, "24299932244005800117978005500793438667981994951685184390218551551204573253089", "puk2"),
	}
}

func getPrivateKey(t *testing.T) *gost3410.PrivateKey {
	key := &gost3410.PrivateKey{
		C:   gost3410.CurveIdGostR34102001CryptoProXchAParamSet(),
		Key: getBigInt(t, "8100551082987309382040692774861374330127499061554316741502830866978492609026", "prk"),
	}

	_ = key.Public()

	println("Key:", hex.EncodeToString(key.Raw()))
	return key
}

func Test_parsePublicKey(t *testing.T) {
	type args struct {
		cert string
	}
	tests := []struct {
		name    string
		args    args
		want    *gost3410.PublicKey
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				cert: validCert,
			},
			want:    getPublicKey(t),
			wantErr: false,
		},
		{
			name: "fail der",
			args: args{
				cert: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail cert",
			args: args{
				cert: badCertOrKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail pub",
			args: args{
				cert: badPub,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePublicKey(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePublicKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHash(t *testing.T) {
	type args struct {
		cert string
	}
	tests := []struct {
		name    string
		args    args
		want    hash.Hash
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				cert: validCert,
			},
			want:    gost34112012256.New(),
			wantErr: false,
		},
		{
			name: "fail der",
			args: args{
				cert: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail cert",
			args: args{
				cert: badCertOrKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail pub",
			args: args{
				cert: badPub,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHash(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePrivateKey(t *testing.T) {
	type args struct {
		key string
		pub *gost3410.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    *gost3410.PrivateKey
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				key: validKey,
				pub: getPublicKey(t),
			},
			want:    getPrivateKey(t),
			wantErr: false,
		},
		{
			name: "fail der",
			args: args{
				key: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail cert",
			args: args{
				key: badCertOrKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail priv",
			args: args{
				key: badKey,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail equal",
			args: args{
				key: validKey,
				pub: getPublicKey2(t),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePrivateKey(tt.args.key, tt.args.pub)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePrivateKey() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestNewGostCrypto(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    sspvo.Crypto
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				opts: []Option{SetCert(validCert), SetKey(validKey)},
			},
			want: &GostCrypto{
				Crypto: Crypto{
					opts: &options{
						cert: validCert,
						key:  validKey,
					},
					hash: gost34112012256.New(),
				},
				privateKey: getPrivateKey(t),
				publicKey:  getPublicKey(t),
			},
			wantErr: false,
		},
		{
			name: "fail crypto",
			args: args{
				opts: []Option{SetKey(validKey)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail pub",
			args: args{
				opts: []Option{SetCert("BAD"), SetKey(validKey)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail priv",
			args: args{
				opts: []Option{SetCert(validCert), SetKey("BAD")},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGostCrypto(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGostCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGostCrypto() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGostCrypto_GetVerifyCrypto(t *testing.T) {
	type fields struct {
		Crypto     Crypto
		privateKey *gost3410.PrivateKey
		publicKey  *gost3410.PublicKey
	}
	type args struct {
		cert string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    sspvo.Crypto
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				cert: validCert,
			},
			want: &GostCrypto{
				Crypto: Crypto{
					opts: &options{
						cert: validCert,
					},
					hash: gost34112012256.New(),
				},
				publicKey: getPublicKey(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GostCrypto{
				Crypto:     tt.fields.Crypto,
				privateKey: tt.fields.privateKey,
				publicKey:  tt.fields.publicKey,
			}
			got, err := c.GetVerifyCrypto(tt.args.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVerifyCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVerifyCrypto() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGostCrypto_Hash(t *testing.T) {
	type fields struct {
		Crypto     Crypto
		privateKey *gost3410.PrivateKey
		publicKey  *gost3410.PublicKey
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
				Crypto: Crypto{
					hash: gost34112012256.New(),
				},
			},
			args: args{
				data: []byte("test"),
			},
			wantHash: testHash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GostCrypto{
				Crypto:     tt.fields.Crypto,
				privateKey: tt.fields.privateKey,
				publicKey:  tt.fields.publicKey,
			}
			gotHash := c.Hash(tt.args.data)
			println("Hash:", hex.EncodeToString(gotHash))
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("Hash() = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}

func TestGostCrypto_Sign(t *testing.T) {
	type fields struct {
		Crypto     Crypto
		privateKey *gost3410.PrivateKey
		publicKey  *gost3410.PublicKey
	}
	type args struct {
		digest []byte
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSign bool
		wantErr  bool
	}{
		{
			name: "ok",
			fields: fields{
				privateKey: getPrivateKey(t),
				publicKey:  getPublicKey(t),
			},
			args: args{
				digest: testHash,
			},
			wantSign: true,
			wantErr:  false,
		},
		{
			name:     "fail private key",
			fields:   fields{},
			args:     args{},
			wantSign: false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GostCrypto{
				Crypto:     tt.fields.Crypto,
				privateKey: tt.fields.privateKey,
				publicKey:  tt.fields.publicKey,
			}
			gotSign, err := c.Sign(tt.args.digest)
			println("Sign:", hex.EncodeToString(gotSign))
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantSign && len(gotSign) == 0 {
				t.Errorf("Sign() gotSign = %v, want %v", gotSign, tt.wantSign)
			}
		})
	}
}

func TestGostCrypto_Verify(t *testing.T) {
	type fields struct {
		Crypto     Crypto
		privateKey *gost3410.PrivateKey
		publicKey  *gost3410.PublicKey
	}
	type args struct {
		sign   []byte
		digest []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				publicKey: getPublicKey(t),
			},
			args: args{
				sign:   []byte{224, 120, 158, 127, 219, 86, 176, 237, 99, 21, 16, 203, 119, 251, 188, 139, 169, 1, 30, 27, 250, 218, 55, 69, 89, 214, 112, 151, 200, 209, 90, 115, 159, 119, 33, 65, 70, 147, 116, 235, 182, 53, 7, 149, 60, 173, 186, 99, 7, 35, 214, 25, 217, 132, 211, 167, 155, 223, 247, 55, 47, 9, 171, 222},
				digest: testHash,
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "false",
			fields: fields{
				publicKey: getPublicKey(t),
			},
			args: args{
				sign:   []byte{225, 120, 158, 127, 219, 86, 176, 237, 99, 21, 16, 203, 119, 251, 188, 139, 169, 1, 30, 27, 250, 218, 55, 69, 89, 214, 112, 151, 200, 209, 90, 115, 159, 119, 33, 65, 70, 147, 116, 235, 182, 53, 7, 149, 60, 173, 186, 99, 7, 35, 214, 25, 217, 132, 211, 167, 155, 223, 247, 55, 47, 9, 171, 222},
				digest: testHash,
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "ok extern",
			fields: fields{
				privateKey: getPrivateKey(t),
				publicKey:  getPublicKey(t),
			},
			args: args{
				sign:   signExtern,
				digest: testHash,
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "ok extern2",
			fields: fields{
				privateKey: getPrivateKey(t),
				publicKey:  getPublicKey(t),
			},
			args: args{
				sign:   signExtern2,
				digest: testHash,
			},
			wantOk:  true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GostCrypto{
				Crypto:     tt.fields.Crypto,
				privateKey: tt.fields.privateKey,
				publicKey:  tt.fields.publicKey,
			}
			gotOk, err := c.Verify(tt.args.sign, tt.args.digest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Verify() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
