package scalingo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"

	httpclient "github.com/Scalingo/go-scalingo/v9/http"
	"github.com/Scalingo/go-scalingo/v9/http/httpmock"
)

func MockAuth(ctrl *gomock.Controller) *httpmock.MockClient {
	mock := httpmock.NewMockClient(ctrl)

	mock.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ *httpclient.APIRequest) (*http.Response, error) {
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			return nil, err
		}

		return &http.Response{
			Body: io.NopCloser(bytes.NewBuffer([]byte(fmt.Sprintf(`{"token": "%v"}`, jwt)))),
		}, nil
	}).AnyTimes()
	return mock
}
