package scalingo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	httpclient "github.com/Scalingo/go-scalingo/http"
	"github.com/Scalingo/go-scalingo/http/httpmock"
	"github.com/dgrijalva/jwt-go"
	gomock "github.com/golang/mock/gomock"
)

func MockAuth(ctrl *gomock.Controller) *httpmock.MockClient {
	mock := httpmock.NewMockClient(ctrl)

	mock.EXPECT().Do(gomock.Any()).DoAndReturn(func(_ *httpclient.APIRequest) (*http.Response, error) {
		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			return nil, err
		}

		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(fmt.Sprintf(`{"token": "%v"}`, jwt)))),
		}, nil
	}).AnyTimes()
	return mock
}
