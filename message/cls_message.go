package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type CLSMessage struct {
	Message
}

func NewCLSMessage(cls string) *CLSMessage {
	msg := &CLSMessage{}
	msg.Init()
	msg.UpdateJWTFields(setCLS(cls))

	return msg
}

func (m *CLSMessage) PathMethod() string {
	return pathMethodCLS
}

func (m *CLSMessage) Response() sspvo.Response {
	return response.NewCLSResponse(m.Fields[sspvo.FieldCLS])
}
