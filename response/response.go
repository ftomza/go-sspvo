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
	skipVerify    bool
	cryptoHandler func(cert string) (sspvo.Crypto, error)
}

func (r *SignResponse) SetCryptoHandler(handler func(cert string) (sspvo.Crypto, error)) {
	r.cryptoHandler = handler
}

func (r *SignResponse) SetSkipVerify(skip bool) {
	r.skipVerify = skip
}

func (r *SignResponse) Data() ([]byte, error) {
	data, err := r.Response.Data()
	if err != nil {
		return nil, err
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
	parts := strings.Split(responseToken, ".")
	if parts[0] == "" {
		return res
	}
	res.Header = parts[0]
	if len(parts) > 2 {
		res.Sign = parts[2]
		res.Payload = parts[1]
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

	err = json.Unmarshal(headerData, &headerStruct)
	if err != nil {
		return false, err
	}

	certData, err := base64.StdEncoding.DecodeString(headerStruct.Cert64)
	if err != nil {
		return false, err
	}

	crypto, err := r.cryptoHandler(string(certData))
	if err != nil {
		return false, err
	}

	data4digest := fmt.Sprintf("%s.%s", token.Header, token.Payload)
	digest := crypto.Hash([]byte(data4digest))

	ok, err := crypto.Verify([]byte(token.Sign), digest)
	if err != nil {
		return false, err
	}

	return ok, nil
}

type CLSResponse struct {
	Response
}

func NewCLSResponse(cls string) *CLSResponse {
	return &CLSResponse{}
}

type ActionResponse struct {
	SignResponse
}

func NewActionResponse(action, dataType string) *ActionResponse {
	return &ActionResponse{}
}

type ConfirmResponse struct {
	SignResponse
}

func NewConfirmResponse() *ConfirmResponse {
	return &ConfirmResponse{}
}

type InfoResponse struct {
	SignResponse
}

func NewInfoResponse() *InfoResponse {
	return &InfoResponse{}
}

type InfoAllResponse struct {
	Response
}

func NewInfoAllResponse() *InfoAllResponse {
	return &InfoAllResponse{}
}
