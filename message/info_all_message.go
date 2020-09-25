package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type InfoAllMessage struct {
	Message
}

func NewInfoAllMessage() *InfoAllMessage {
	msg := &InfoAllMessage{}
	msg.Init()

	return msg
}

func (m *InfoAllMessage) PathMethod() string {
	return pathMethodInfo
}

func (m *InfoAllMessage) Response() sspvo.Response {
	return response.NewInfoAllResponse()
}
