package scalingo

import (
	"bytes"
	"io/ioutil"
	"net/http"

	httpclient "github.com/Scalingo/go-scalingo/http"
	"github.com/Scalingo/go-scalingo/http/httpmock"
	gomock "github.com/golang/mock/gomock"
)

func MockAuth(ctrl *gomock.Controller) *httpmock.MockClient {
	mock := httpmock.NewMockClient(ctrl)
	mock.EXPECT().Do(gomock.Any()).DoAndReturn(func(_ *httpclient.APIRequest) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"token": "toto"}`))),
		}, nil
	}).AnyTimes()
	return mock
}
