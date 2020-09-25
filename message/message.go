package message

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ftomza/go-sspvo"
)

const (
	pathMethodAction  = "token/new"
	pathMethodCLS     = "cls/request"
	pathMethodInfo    = "token/info"
	pathMethodConfirm = "token/confirm"
)

func SetOGRN(ogrn string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldOGRN] = ogrn
	}
}

func SetKPP(kpp string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldKPP] = kpp
	}
}

func setCert(cert string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldCert] = base64.StdEncoding.EncodeToString([]byte(cert))
	}
}

func setToken(token *sspvo.Token) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldToken] = fmt.Sprintf("%s.%s.%s", token.Header, token.Payload, token.Sign)
	}
}

func setIdJWT(idJWT string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldIdJWT] = idJWT
	}
}

func setCLS(cls string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldCLS] = cls
	}
}

func setAction(action string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldAction] = action
	}
}

func setDataType(dataType string) sspvo.Fields {
	return func(f sspvo.JWTFields) {
		f[sspvo.FieldDataType] = dataType
	}
}

type Message struct {
	Fields sspvo.JWTFields
}

func (m *Message) UpdateJWTFields(fields ...sspvo.Fields) sspvo.Message {
	for _, field := range fields {
		field(m.Fields)
	}
	return m
}

func (m *Message) Init() {
	m.Fields = sspvo.JWTFields{}
}

func (m *Message) PathMethod() string {
	panic("not implement")
}

func (m *Message) GetJWT() ([]byte, error) {
	return json.Marshal(m.Fields)
}

func (m *Message) Response() sspvo.Response {

	panic("not implement")
}

type SignMessage struct {
	Message
	crypto sspvo.Crypto
	data   []byte
}

func (m *SignMessage) GetJWT() ([]byte, error) {
	token, err := m.sign()
	if err != nil {
		return nil, fmt.Errorf("SignMessage: %w", err)
	}
	m.UpdateJWTFields(setToken(token))
	return m.Message.GetJWT()
}

func (m *SignMessage) Init(crypto sspvo.Crypto, data []byte) {
	m.Message.Init()
	m.crypto = crypto
	m.data = data
}

func (m *SignMessage) sign() (*sspvo.Token, error) {
	m.UpdateJWTFields(setCert(m.crypto.GetCert()))

	jsonHeader, err := json.Marshal(m.Fields)
	if err != nil {
		return nil, err
	}
	header := base64.StdEncoding.EncodeToString(jsonHeader)

	payload := ""
	if m.data != nil {
		payload = base64.StdEncoding.EncodeToString(m.data)
	}

	data4sign := []byte(fmt.Sprintf("%s.%s", header, payload))
	digest := m.crypto.Hash(data4sign)
	signDigest, err := m.crypto.Sign(digest)
	if err != nil {
		return nil, err
	}
	if ok, err := m.crypto.Verify(signDigest, digest); err != nil {
		return nil, err
	} else if !ok {
		return nil, sspvo.ErrBadSign
	}
	sign := base64.StdEncoding.EncodeToString(signDigest)

	return &sspvo.Token{
		Header:  header,
		Payload: payload,
		Sign:    sign,
	}, nil
}
