package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type InfoMessage struct {
	SignMessage
}

func NewInfoMessage(crypto sspvo.Crypto, idJWT string) *InfoMessage {
	msg := &InfoMessage{}
	msg.Init(crypto, nil)
	msg.UpdateJWTFields(setIdJWT(idJWT))

	return msg
}

func (m *InfoMessage) PathMethod() string {
	return pathMethodInfo
}

func (m *InfoMessage) Response() sspvo.Response {
	resp := response.NewInfoResponse()
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
