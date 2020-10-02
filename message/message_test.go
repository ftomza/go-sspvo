/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/ftomza/go-sspvo"

	"github.com/ftomza/go-sspvo/mocks"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	client *mocks.Client
	crypto *mocks.Crypto
}

func (suite *MessageTestSuite) SetupTest() {
	suite.client = new(mocks.Client)
	suite.crypto = new(mocks.Crypto)
}

func Test_EntUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

func (suite *MessageTestSuite) TearDownTest() {
}

func (suite *MessageTestSuite) TestMessage_GetJWT() {
	suite.Run("ok message", func() {
		msg := Message{
			Fields: sspvo.JWTFields{sspvo.FieldAction: "Add"},
		}
		b, err := msg.GetJWT()
		suite.NoError(err)
		suite.Equal([]byte("{\"action\":\"Add\"}"), b)
	})
	suite.Run("ok sign message", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(true, nil).Once()
		suite.crypto.On("Sign", mock.Anything).Return([]byte("TEST"), nil).Once()
		msg := func() sspvo.Message {
			return &SignMessage{
				Message: Message{
					Fields: sspvo.JWTFields{},
				},
				crypto: suite.crypto,
			}
		}
		b, err := msg().GetJWT()
		suite.NoError(err)
		suite.Equal([]byte("{\"token\":\"eyJDZXJ0NjQiOiJURVNUIn0=..VEVTVA==\"}"), b)
	})
	suite.Run("fail sign", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Sign", mock.Anything).Return(nil, errors.New("fail")).Once()
		msg := func() sspvo.Message {
			return &SignMessage{
				Message: Message{
					Fields: sspvo.JWTFields{},
				},
				crypto: suite.crypto,
			}
		}
		b, err := msg().GetJWT()
		suite.Error(err)
		suite.Nil(b)
	})
}

func (suite *MessageTestSuite) TestMessage_signToken() {
	suite.Run("ok", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Sign", mock.Anything).Return([]byte("TEST"), nil).Once()
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(true, nil).Once()
		msg := &SignMessage{
			Message: Message{
				Fields: sspvo.JWTFields{},
			},
			crypto: suite.crypto,
		}
		token, err := msg.signToken()
		suite.NoError(err)
		suite.Equal(&sspvo.Token{
			Header:  "eyJDZXJ0NjQiOiJURVNUIn0=",
			Payload: "",
			Sign:    "VEVTVA==",
		}, token)

		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("fail sign", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Sign", mock.Anything).Return(nil, errors.New("fail")).Once()
		msg := &SignMessage{
			Message: Message{
				Fields: sspvo.JWTFields{"BAD": suite},
			},
			crypto: suite.crypto,
		}
		token, err := msg.signToken()
		suite.Error(err)
		suite.Nil(token)

		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("fail verify", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Sign", mock.Anything).Return([]byte("TEST"), nil).Once()
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(false, errors.New("fail")).Once()
		msg := &SignMessage{
			Message: Message{
				Fields: sspvo.JWTFields{"BAD": suite},
			},
			crypto: suite.crypto,
		}
		token, err := msg.signToken()
		suite.Error(err)
		suite.Nil(token)

		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("bad sign", func() {
		suite.crypto.On("GetCert").Return("TEST").Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		suite.crypto.On("Sign", mock.Anything).Return([]byte("TEST"), nil).Once()
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(false, nil).Once()
		msg := &SignMessage{
			Message: Message{
				Fields: sspvo.JWTFields{"BAD": suite},
			},
			crypto: suite.crypto,
		}
		token, err := msg.signToken()
		suite.Error(err)
		suite.Nil(token)

		suite.crypto.AssertExpectations(suite.T())
	})
}

func (suite *MessageTestSuite) Test_setCert() {
	suite.Run("ok", func() {

		fields := sspvo.JWTFields{}
		setCert("VEVTVA==")(fields)
		suite.Equal(fields[sspvo.FieldCert], "VEVTVA==")
	})
	suite.Run("ok pem", func() {

		fields := sspvo.JWTFields{}
		cert := "-----BEGIN CERTIFICATE-----\nVEVTVA==\n-----END CERTIFICATE-----"
		setCert(cert)(fields)
		suite.Equal(fields[sspvo.FieldCert], "VEVTVA==")
	})
}
