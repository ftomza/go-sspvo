package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type ConfirmMessage struct {
	SignMessage
}

func NewConfirmMessage(crypto sspvo.Crypto, idJWT string) *ConfirmMessage {
	msg := &ConfirmMessage{}
	msg.Init(crypto, nil)
	msg.UpdateJWTFields(setIdJWT(idJWT))

	return msg
}

func (m *ConfirmMessage) PathMethod() string {
	return pathMethodConfirm
}

func (m *ConfirmMessage) Response() sspvo.Response {
	resp := response.NewConfirmResponse()
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
