/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package response

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/ftomza/go-sspvo/crypto"

	"github.com/ftomza/go-sspvo/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ftomza/go-sspvo"
)

var (
	validResponseToken = `
eyJJREpXVCI6MTM4NjE3MSwiZGF0YV90eXBlIjoiRXJyb3IiLCJDZXJ0NjQiOiJNSUlFWVRDQ0JCQ2dBd0lCQWdJVEVnQkJobTI1UjhJVkc2RnlNZ0FCQUVHR2JUQUlCZ1lxaFFNQ0FnTXdmekVqTUNFR0NTcUdTSWIzRFFFSkFSWVVjM1Z3Y0c5eWRFQmpjbmx3ZEc5d2NtOHVjblV4Q3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRSEV3Wk5iM05qYjNjeEZ6QVZCZ05WQkFvVERrTlNXVkJVVHkxUVVrOGdURXhETVNFd0h3WURWUVFERXhoRFVsbFFWRTh0VUZKUElGUmxjM1FnUTJWdWRHVnlJREl3SGhjTk1qQXdNakk0TVRNek1qSTBXaGNOTWpBd05USTRNVE0wTWpJMFdqQ0NBV2d4R0RBV0JnVXFoUU5rQVJJTk5UQTROemMwTmpZNU56VXlPREU2TURnR0ExVUVDUXd4MEtQUXV5N1FuOUdBMExYUmdkQzkwTFhRdmRHQjBMclF1TkM1SU5DUzBMRFF1eURRdEM0eE9TRFJnZEdDMFlBdU1URVRNQkVHQTFVRUNnd0swS0hRc05DODBZdlF1VEVWTUJNR0ExVUVCd3dNMEp6UXZ0R0IwTHJRc3RDd01Sb3dHQVlJS29VREE0RURBUUVTRERBd056Z3hORFV3T0RreU1URWJNQmtHQ1NxR1NJYjNEUUVKQVJZTWJXRnBiRUJ0WVdsc0xuSjFNUm93R0FZSUtvVURBNEVEQVFFU0RERXhNVEV4TVRFeE1URXhNVEVZTUJZR0JTcUZBMlFCRWcweE1URXhNVEV4TVRFeE1URXhNUXN3Q1FZRFZRUUdFd0pTVlRFaU1DQUdBMVVFS2d3WjBKalFzdEN3MEwwZzBKalFzdEN3MEwzUXZ0Q3kwTGpSaHpFUk1BOEdBMVVFREF3STBLRFJnOUM3MFl3eEdqQVlCZ05WQkFNTUVkQ1kwS0VnMEtIUXNOQzgwTDdRczlDK01SVXdFd1lEVlFRRURBelFtTkN5MExEUXZkQyswTEl3WmpBZkJnZ3FoUU1IQVFFQkFUQVRCZ2NxaFFNQ0FpUUFCZ2dxaFFNSEFRRUNBZ05EQUFSQTJNV2xNREpBUXBYcjJyU3hGQWk4UGh0UUozNmJ3NHF3YVNWcTkwNTNZektQVGJoMXdVNXdGMzJZM2lFUW5rUWwzU1hiWTQzWnhXZUVQZmR2UHZhSEw2T0NBWE13Z2dGdk1CTUdBMVVkSlFRTU1Bb0dDQ3NHQVFVRkJ3TUVNQXNHQTFVZER3UUVBd0lFOERBZEJnTlZIUTRFRmdRVTF4SmxuM2hVZUEvZlJONnlUeWs3T1Q2RFNVMHdId1lEVlIwakJCZ3dGb0FVVG9NK0ZHbnY3RjE2bFN0ZkVmNDNNaFpKVlNzd1hBWURWUjBmQkZVd1V6QlJvRStnVFlaTGFIUjBjRG92TDNSbGMzUmpZUzVqY25sd2RHOXdjbTh1Y25VdlEyVnlkRVZ1Y205c2JDOURVbGxRVkU4dFVGSlBKVEl3VkdWemRDVXlNRU5sYm5SbGNpVXlNRElvTVNrdVkzSnNNSUdzQmdnckJnRUZCUWNCQVFTQm56Q0JuREJrQmdnckJnRUZCUWN3QW9aWWFIUjBjRG92TDNSbGMzUmpZUzVqY25sd2RHOXdjbTh1Y25VdlEyVnlkRVZ1Y205c2JDOTBaWE4wTFdOaExUSXdNVFJmUTFKWlVGUlBMVkJTVHlVeU1GUmxjM1FsTWpCRFpXNTBaWElsTWpBeUtERXBMbU55ZERBMEJnZ3JCZ0VGQlFjd0FZWW9hSFIwY0RvdkwzUmxjM1JqWVM1amNubHdkRzl3Y204dWNuVXZiMk56Y0M5dlkzTndMbk55WmpBSUJnWXFoUU1DQWdNRFFRQnJzdk1TYkk4enJjR0FUcXBwVFJjRkVCeEI5UWVUODBZb1orWm13ZUVaazN0TFhGdk1RN2V6SG9nc0tXMHNMNlhEVnpKUUt6RG1ISWZPNDlOYlVtN1kifQ==.PFBhY2thZ2VEYXRhPjxFbnRpdHlUeXBlPmVkaXRfYXBwbGljYXRpb25fc3RhdHVzPC9FbnRpdHlUeXBlPjxFcnJvcj48RXJyb3JDb2RlPjQwNTY8L0Vycm9yQ29kZT48RXJyb3JEZXNjcmlwdGlvbj7Ql9Cw0Y/QstC70LXQvdC40LUg0YEgVUlERXBndSAmIzM0Ozg2NTYzMzQ4NCYjMzQ7INC90LUg0L3QsNC50LTQtdC90L4uPC9FcnJvckRlc2NyaXB0aW9uPjxVSURFcGd1Pjg2NTYzMzQ4NDwvVUlERXBndT48L0Vycm9yPjwvUGFja2FnZURhdGE+.4okpt8vyBjP7KdRNTiKLPgex+Fk3kjh2xmNJsutzMZu6KrGxmxawEfwM9Uaj1Ht7U2hOWPMGOchdrLOexR0khA==
`
	badResponseToken = `
eyJJREpXVCI6MTM4NjE3MSwiZGF0YV90eXBlIjoiRXJyb3IiLCJDZXJ0NjQiOiJNSUlFWVRDQ0JCQ2dBd0lCQWdJVEVnQkJobTI1UjhJVkc2RnlNZ0FCQUVHR2JUQUlCZ1lxaFFNQ0FnTXdmekVqTUNFR0NTcUdTSWIzRFFFSkFSWVVjM1Z3Y0c5eWRFQmpjbmx3ZEc5d2NtOHVjblV4Q3pBSkJnTlZCQVlUQWxKVk1ROHdEUVlEVlFRSEV3Wk5iM05qYjNjeEZ6QVZCZ05WQkFvVERrTlNXVkJVVHkxUVVrOGdURXhETVNFd0h3WURWUVFERXhoRFVsbFFWRTh0VUZKUElGUmxjM1FnUTJWdWRHVnlJREl3SGhjTk1qQXdNakk0TVRNek1qSTBXaGNOTWpBd05USTRNVE0wTWpJMFdqQ0NBV2d4R0RBV0JnVXFoUU5rQVJJTk5UQTROemMwTmpZNU56VXlPREU2TURnR0ExVUVDUXd4MEtQUXV5N1FuOUdBMExYUmdkQzkwTFhRdmRHQjBMclF1TkM1SU5DUzBMRFF1eURRdEM0eE9TRFJnZEdDMFlBdU1URVRNQkVHQTFVRUNnd0swS0hRc05DODBZdlF1VEVWTUJNR0ExVUVCd3dNMEp6UXZ0R0IwTHJRc3RDd01Sb3dHQVlJS29VREE0RURBUUVTRERBd056Z3hORFV3T0RreU1URWJNQmtHQ1NxR1NJYjNEUUVKQVJZTWJXRnBiRUJ0WVdsc0xuSjFNUm93R0FZSUtvVURBNEVEQVFFU0RERXhNVEV4TVRFeE1URXhNVEVZTUJZR0JTcUZBMlFCRWcweE1URXhNVEV4TVRFeE1URXhNUXN3Q1FZRFZRUUdFd0pTVlRFaU1DQUdBMVVFS2d3WjBKalFzdEN3MEwwZzBKalFzdEN3MEwzUXZ0Q3kwTGpSaHpFUk1BOEdBMVVFREF3STBLRFJnOUM3MFl3eEdqQVlCZ05WQkFNTUVkQ1kwS0VnMEtIUXNOQzgwTDdRczlDK01SVXdFd1lEVlFRRURBelFtTkN5MExEUXZkQyswTEl3WmpBZkJnZ3FoUU1IQVFFQkFUQVRCZ2NxaFFNQ0FpUUFCZ2dxaFFNSEFRRUNBZ05EQUFSQTJNV2xNREpBUXBYcjJyU3hGQWk4UGh0UUozNmJ3NHF3YVNWcTkwNTNZektQVGJoMXdVNXdGMzJZM2lFUW5rUWwzU1hiWTQzWnhXZUVQZmR2UHZhSEw2T0NBWE13Z2dGdk1CTUdBMVVkSlFRTU1Bb0dDQ3NHQVFVRkJ3TUVNQXNHQTFVZER3UUVBd0lFOERBZEJnTlZIUTRFRmdRVTF4SmxuM2hVZUEvZlJONnlUeWs3T1Q2RFNVMHdId1lEVlIwakJCZ3dGb0FVVG9NK0ZHbnY3RjE2bFN0ZkVmNDNNaFpKVlNzd1hBWURWUjBmQkZVd1V6QlJvRStnVFlaTGFIUjBjRG92TDNSbGMzUmpZUzVqY25sd2RHOXdjbTh1Y25VdlEyVnlkRVZ1Y205c2JDOURVbGxRVkU4dFVGSlBKVEl3VkdWemRDVXlNRU5sYm5SbGNpVXlNRElvTVNrdVkzSnNNSUdzQmdnckJnRUZCUWNCQVFTQm56Q0JuREJrQmdnckJnRUZCUWN3QW9aWWFIUjBjRG92TDNSbGMzUmpZUzVqY25sd2RHOXdjbTh1Y25VdlEyVnlkRVZ1Y205c2JDOTBaWE4wTFdOaExUSXdNVFJmUTFKWlVGUlBMVkJTVHlVeU1GUmxjM1FsTWpCRFpXNTBaWElsTWpBeUtERXBMbU55ZERBMEJnZ3JCZ0VGQlFjd0FZWW9hSFIwY0RvdkwzUmxjM1JqWVM1amNubHdkRzl3Y204dWNuVXZiMk56Y0M5dlkzTndMbk55WmpBSUJnWXFoUU1DQWdNRFFRQnJzdk1TYkk4enJjR0FUcXBwVFJjRkVCeEI5UWVUODBZb1orWm13ZUVaazN0TFhGdk1RN2V6SG9nc0tXMHNMNlhEVnpKUUt6RG1ISWZPNDlOYlVtN1kifQ==.VEVTVA==.VEVTVA==
`
	falseResponseToken = `
eyJJREpXVCI6MTEwMjI5NiwiZGF0YV90eXBlIjoiU3VjY2VzcyIsIkNlcnQ2NCI6Ik1JSUVZVENDQkJDZ0F3SUJBZ0lURWdCRlVyT2NyUnZsQUx6NG1RQUJBRVZTc3pBSUJnWXFoUU1DQWdNd2Z6RWpNQ0VHQ1NxR1NJYjNEUUVKQVJZVWMzVndjRzl5ZEVCamNubHdkRzl3Y204dWNuVXhDekFKQmdOVkJBWVRBbEpWTVE4d0RRWURWUVFIRXdaTmIzTmpiM2N4RnpBVkJnTlZCQW9URGtOU1dWQlVUeTFRVWs4Z1RFeERNU0V3SHdZRFZRUURFeGhEVWxsUVZFOHRVRkpQSUZSbGMzUWdRMlZ1ZEdWeUlESXdIaGNOTWpBd05USTRNVFF4TkRJeFdoY05NakF3T0RJNE1UUXlOREl4V2pDQ0FXZ3hHREFXQmdVcWhRTmtBUklOTlRBNE56YzBOalk1TnpVeU9ERTZNRGdHQTFVRUNRd3gwS1BRdXk3UW45R0EwTFhSZ2RDOTBMWFF2ZEdCMExyUXVOQzVJTkNTMExEUXV5RFF0QzR4T1NEUmdkR0MwWUF1TVRFVE1CRUdBMVVFQ2d3SzBLSFFzTkM4MFl2UXVURVZNQk1HQTFVRUJ3d00wSnpRdnRHQjBMclFzdEN3TVJvd0dBWUlLb1VEQTRFREFRRVNEREF3TnpneE5EVXdPRGt5TVRFYk1Ca0dDU3FHU0liM0RRRUpBUllNYldGcGJFQnRZV2xzTG5KMU1Sb3dHQVlJS29VREE0RURBUUVTRERFeE1URXhNVEV4TVRFeE1URVlNQllHQlNxRkEyUUJFZzB4TVRFeE1URXhNVEV4TVRFeE1Rc3dDUVlEVlFRR0V3SlNWVEVpTUNBR0ExVUVLZ3daMEpqUXN0Q3cwTDBnMEpqUXN0Q3cwTDNRdnRDeTBMalJoekVSTUE4R0ExVUVEQXdJMEtEUmc5QzcwWXd4R2pBWUJnTlZCQU1NRWRDWTBLRWcwS0hRc05DODBMN1FzOUMrTVJVd0V3WURWUVFFREF6UW1OQ3kwTERRdmRDKzBMSXdaakFmQmdncWhRTUhBUUVCQVRBVEJnY3FoUU1DQWlRQUJnZ3FoUU1IQVFFQ0FnTkRBQVJBMk1XbE1ESkFRcFhyMnJTeEZBaThQaHRRSjM2Ync0cXdhU1ZxOTA1M1l6S1BUYmgxd1U1d0YzMlkzaUVRbmtRbDNTWGJZNDNaeFdlRVBmZHZQdmFITDZPQ0FYTXdnZ0Z2TUJNR0ExVWRKUVFNTUFvR0NDc0dBUVVGQndNRU1Bc0dBMVVkRHdRRUF3SUU4REFkQmdOVkhRNEVGZ1FVMXhKbG4zaFVlQS9mUk42eVR5azdPVDZEU1Uwd0h3WURWUjBqQkJnd0ZvQVVUb00rRkdudjdGMTZsU3RmRWY0M01oWkpWU3N3WEFZRFZSMGZCRlV3VXpCUm9FK2dUWVpMYUhSMGNEb3ZMM1JsYzNSallTNWpjbmx3ZEc5d2NtOHVjblV2UTJWeWRFVnVjbTlzYkM5RFVsbFFWRTh0VUZKUEpUSXdWR1Z6ZENVeU1FTmxiblJsY2lVeU1ESW9NU2t1WTNKc01JR3NCZ2dyQmdFRkJRY0JBUVNCbnpDQm5EQmtCZ2dyQmdFRkJRY3dBb1pZYUhSMGNEb3ZMM1JsYzNSallTNWpjbmx3ZEc5d2NtOHVjblV2UTJWeWRFVnVjbTlzYkM5MFpYTjBMV05oTFRJd01UUmZRMUpaVUZSUExWQlNUeVV5TUZSbGMzUWxNakJEWlc1MFpYSWxNakF5S0RFcExtTnlkREEwQmdnckJnRUZCUWN3QVlZb2FIUjBjRG92TDNSbGMzUmpZUzVqY25sd2RHOXdjbTh1Y25VdmIyTnpjQzl2WTNOd0xuTnlaakFJQmdZcWhRTUNBZ01EUVFCa29VZnVjRy8xUnZUQ09ac2RVc1hvcWZjQW9Hd2VFbWtaMFJESXNkZVlpaHF0TFQ1cERVMW91cVZvK2dRcmt5NFprQ2dPdnpIRXFXeUcwYWJyc2d6YyJ9.PFBhY2thZ2VEYXRhPjxUWVBFPmFjaGlldmVtZW50czwvVFlQRT48TUVTU0FHRT7QrdC70LXQvNC10L3RgiDQstGL0LPRgNGD0LbQtdC9PC9NRVNTQUdFPjxCQVNFNjRGSUxFPmFHa0s8L0JBU0U2NEZJTEU+PEVYVD5wZGY8L0VYVD48L1BhY2thZ2VEYXRhPg==.vVlB9zxuUVDaY6tQogEjOoYnoxDNpCf7Ah5bgDrpm94ii7nXHZ2zrCxmnWkNZITlR7ZmtZI57SIpqaxS2UanNg==
`
	dataToken      = []byte(`{"ResponseToken": "` + strings.TrimSpace(validResponseToken) + `"}`)
	dataEmptyToken = []byte(`{"ResponseToken": ""}`)
)

type ResponseTestSuite struct {
	suite.Suite
	crypto *mocks.Crypto
}

func (suite *ResponseTestSuite) SetupTest() {
	suite.crypto = new(mocks.Crypto)
}

func Test_EntUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ResponseTestSuite))
}

func (suite *ResponseTestSuite) TearDownTest() {
}

func (suite *ResponseTestSuite) TestSignResponse_Data() {
	suite.Run("ok", func() {
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(true, nil).Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()

		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: dataToken,
				},
			},
			skipVerify: false,
			cryptoHandler: func(cert string) (sspvo.Crypto, error) {
				return suite.crypto, nil
			},
		}

		b, err := response.Data()
		suite.NoError(err)
		suite.Equal(dataToken, b)

		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("ok empty token", func() {
		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: dataEmptyToken,
				},
			},
			skipVerify: false,
			cryptoHandler: func(cert string) (sspvo.Crypto, error) {
				return suite.crypto, nil
			},
		}

		b, err := response.Data()
		suite.NoError(err)
		suite.Equal(dataEmptyToken, b)

		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("fail response", func() {

		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusBadRequest,
					Body: []byte("BAD"),
				},
			},
		}

		b, err := response.Data()
		suite.Error(err)
		suite.Equal([]byte("BAD"), b)
		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("fail token", func() {

		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: []byte("BAD"),
				},
			},
		}

		b, err := response.Data()
		suite.Error(err)
		suite.Nil(b)
		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("fail verify", func() {
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(false, errors.New("fail")).Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: dataToken,
				},
			},
			cryptoHandler: func(cert string) (sspvo.Crypto, error) {
				return suite.crypto, nil
			},
		}

		b, err := response.Data()
		suite.Error(err)
		suite.Nil(b)
		suite.crypto.AssertExpectations(suite.T())
	})
	suite.Run("false verify", func() {
		suite.crypto.On("Verify", mock.Anything, mock.Anything).Return(false, nil).Once()
		suite.crypto.On("Hash", mock.Anything).Return([]byte("TEST")).Once()
		response := &SignResponse{
			Response: Response{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: dataToken,
				},
			},
			cryptoHandler: func(cert string) (sspvo.Crypto, error) {
				return suite.crypto, nil
			},
		}

		b, err := response.Data()
		suite.Error(err)
		suite.Nil(b)
		suite.crypto.AssertExpectations(suite.T())
	})
}

func TestResponse_Data(t *testing.T) {
	type fields struct {
		resp *sspvo.ClientResponse
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				resp: &sspvo.ClientResponse{
					Code: http.StatusOK,
					Body: []byte("test"),
				},
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "fail code",
			fields: fields{
				resp: &sspvo.ClientResponse{
					Code: http.StatusUnauthorized,
					Body: []byte("test"),
				},
			},
			want:    []byte("test"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				resp: tt.fields.resp,
				err:  tt.fields.err,
			}
			got, err := r.Data()
			if (err != nil) != tt.wantErr {
				t.Errorf("Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignResponse_parseResponseToken(t *testing.T) {
	type fields struct {
		Response      Response
		skipVerify    bool
		cryptoHandler func(cert string) (sspvo.Crypto, error)
	}
	type args struct {
		responseToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *sspvo.Token
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				"ONE.TWO.NINE",
			},
			wantRes: &sspvo.Token{
				Header:  "ONE",
				Payload: "TWO",
				Sign:    "NINE",
			},
		},
		{
			name:   "ok two",
			fields: fields{},
			args: args{
				"ONE.NINE",
			},
			wantRes: &sspvo.Token{
				Header: "ONE",
				Sign:   "NINE",
			},
		},
		{
			name:   "ok none",
			fields: fields{},
			args: args{
				"..",
			},
			wantRes: &sspvo.Token{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SignResponse{
				Response:      tt.fields.Response,
				skipVerify:    tt.fields.skipVerify,
				cryptoHandler: tt.fields.cryptoHandler,
			}
			if gotRes := r.parseResponseToken(tt.args.responseToken); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("parseResponseToken() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestSignResponse_verifyResponseToken(t *testing.T) {
	type fields struct {
		Response      Response
		skipVerify    bool
		cryptoHandler func(cert string) (sspvo.Crypto, error)
	}
	type args struct {
		responseToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
				cryptoHandler: func(cert string) (sspvo.Crypto, error) {
					return crypto.NewGostCrypto(crypto.SetCert(cert))
				},
			},
			args: args{
				responseToken: validResponseToken,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "false",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
				cryptoHandler: func(cert string) (sspvo.Crypto, error) {
					return crypto.NewGostCrypto(crypto.SetCert(cert))
				},
			},
			args: args{
				responseToken: falseResponseToken,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ok skip",
			fields: fields{
				Response:   Response{},
				skipVerify: true,
			},
			args: args{
				responseToken: "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ok header empty",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
			},
			args: args{
				responseToken: "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "fail header",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
			},
			args: args{
				responseToken: "BAD.BAD",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail sign",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
			},
			args: args{
				responseToken: "VEVTVA==.BAD.BAD",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail data cert",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
			},
			args: args{
				responseToken: "VEVTVA==.VEVTVA==.VEVTVA==",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail crypto handler",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
				cryptoHandler: func(cert string) (sspvo.Crypto, error) {
					return nil, errors.New("fail")
				},
			},
			args: args{
				responseToken: "e30=.VEVTVA==.VEVTVA==",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail verify",
			fields: fields{
				Response:   Response{},
				skipVerify: false,
				cryptoHandler: func(cert string) (sspvo.Crypto, error) {
					return crypto.NewGostCrypto(crypto.SetCert(cert))
				},
			},
			args: args{
				responseToken: badResponseToken,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SignResponse{
				Response:      tt.fields.Response,
				skipVerify:    tt.fields.skipVerify,
				cryptoHandler: tt.fields.cryptoHandler,
			}
			got, err := r.verifyResponseToken(tt.args.responseToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyResponseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyResponseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
