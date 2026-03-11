package scalingo

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-scalingo/v11/http/httpmock"
)

func MockAuth(ctrl *gomock.Controller) *httpmock.MockClient {
	mock := httpmock.NewMockClient(ctrl)

	mock.EXPECT().DoRequest(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ *httpclient.APIRequest, data any) error {
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		jwt, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			return err
		}

		if data != nil {
			res, _ := data.(*BearerTokenRes)
			res.Token = jwt
		}
		return nil
	}).AnyTimes()
	return mock
}
