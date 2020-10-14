/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package response

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ftomza/go-sspvo"
)

type Response struct {
	resp *sspvo.ClientResponse
	err  error
}

func (r *Response) SetClientResponse(resp *sspvo.ClientResponse) {
	r.resp = resp
}

func (r *Response) ClientResponse() *sspvo.ClientResponse {
	return r.resp
}

func (r *Response) SetError(err error) {
	if err == nil {
		return
	}
	r.err = err
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) Data() ([]byte, error) {
	if r.Error() != nil {
		return nil, r.Error()
	}

	if r.resp.Code > 299 {
		return r.resp.Body, fmt.Errorf("wrong response code: %d", r.resp.Code)
	}

	return r.resp.Body, nil
}

type SignResponse struct {
	Response
	skipVerify bool
	crypto     sspvo.Crypto
}

func (r *SignResponse) Data() ([]byte, error) {
	data, err := r.Response.Data()
	if err != nil {
		return data, err
	}
	signStruct := struct {
		ResponseToken string `json:"ResponseToken"`
	}{}

	err = json.Unmarshal(data, &signStruct)
	if err != nil {
		return nil, fmt.Errorf("SignResponse: %w", err)
	}

	if signStruct.ResponseToken == "" {
		return data, nil
	}

	ok, err := r.verifyResponseToken(signStruct.ResponseToken)
	if err != nil {
		return nil, fmt.Errorf("SignResponse: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("SignResponse: %w", sspvo.ErrBadSign)
	}
	return data, nil
}

func (r *SignResponse) parseResponseToken(responseToken string) (res *sspvo.Token) {
	res = &sspvo.Token{}
	parts := strings.Split(strings.TrimSpace(responseToken), ".")
	if parts[0] == "" {
		return res
	}
	res.Header = parts[0]
	if len(parts) > 2 {
		res.Payload = parts[1]
		res.Sign = parts[2]
	} else {
		res.Sign = parts[1]
	}
	return
}

func (r *SignResponse) verifyResponseToken(responseToken string) (bool, error) {

	if r.skipVerify {
		return true, nil
	}

	token := r.parseResponseToken(responseToken)
	if token.Header == "" {
		return true, nil
	}

	headerStruct := struct {
		Cert64 string `json:"Cert64"`
	}{}

	headerData, err := base64.StdEncoding.DecodeString(token.Header)
	if err != nil {
		return false, err
	}

	sign, err := base64.StdEncoding.DecodeString(token.Sign)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(headerData, &headerStruct)
	if err != nil {
		return false, err
	}

	cert := "-----BEGIN CERTIFICATE-----\n" +
		headerStruct.Cert64 +
		"\n-----END CERTIFICATE-----"

	crypto, err := r.crypto.GetVerifyCrypto(cert)
	if err != nil {
		return false, err
	}

	data4digest := fmt.Sprintf("%s.%s", token.Header, token.Payload)
	digest := crypto.Hash([]byte(data4digest))
	ok, err := crypto.Verify(sign, digest)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func NewResponse() *Response {
	return &Response{}
}

func NewSignResponse(crypto sspvo.Crypto, skipVerify bool) *SignResponse {
	return &SignResponse{
		skipVerify: skipVerify,
		crypto:     crypto,
	}
}
