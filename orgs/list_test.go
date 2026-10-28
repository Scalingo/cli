package orgs

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/scalingomock"
)

type fakeRenderer struct {
	data      []scalingo.Organization
	renderErr error
	rendered  bool
}

func (r *fakeRenderer) SetData(_ context.Context, data []scalingo.Organization) {
	r.data = data
}

func (r *fakeRenderer) Render(_ context.Context) error {
	r.rendered = true
	return r.renderErr
}

func TestList(t *testing.T) {
	organizations := []scalingo.Organization{
		{ID: "org-1", Slug: "acme", Role: "admin"},
		{ID: "org-2", Slug: "globex", Role: "member"},
	}

	tests := map[string]struct {
		listResult       []scalingo.Organization
		listErr          error
		renderErr        error
		expectedError    string
		expectedRendered bool
	}{
		"it should render the organizations": {
			listResult:       organizations,
			expectedRendered: true,
		},
		"it should render an empty list": {
			listResult:       []scalingo.Organization{},
			expectedRendered: true,
		},
		"it should return an error without rendering when listing fails": {
			listErr:       stderrors.New("unauthorized"),
			expectedError: "list organizations",
		},
		"it should return an error when rendering fails": {
			listResult:       organizations,
			renderErr:        stderrors.New("broken pipe"),
			expectedError:    "render organizations list",
			expectedRendered: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			service := scalingomock.NewMockOrganizationsService(gomock.NewController(t))
			service.EXPECT().OrganizationsList(ctx).Return(test.listResult, test.listErr)
			r := &fakeRenderer{renderErr: test.renderErr}

			err := list(ctx, service, r)

			if test.expectedError != "" {
				require.ErrorContains(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.listResult, r.data)
			}
			assert.Equal(t, test.expectedRendered, r.rendered)
		})
	}
}
