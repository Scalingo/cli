package scalingo

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v6/http"
)

type DatabaseTypeVersionPlugin struct {
	ID          string `json:"id"`
	FeatureName string `json:"feature_name"`
	InstallName string `json:"install_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type DatabaseTypeVersion struct {
	ID             string                      `json:"id"`
	DatabaseTypeID string                      `json:"database_type_id"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	Features       []string                    `json:"features"`
	NextUpgrade    *DatabaseTypeVersion        `json:"next_upgrade"`
	AllowedPlugins []DatabaseTypeVersionPlugin `json:"allowed_plugins"`
	Major          int                         `json:"major"`
	Minor          int                         `json:"minor"`
	Patch          int                         `json:"patch"`
	Build          int                         `json:"build"`
}

func (v DatabaseTypeVersion) String() string {
	return fmt.Sprintf("%v.%v.%v-%v", v.Major, v.Minor, v.Patch, v.Build)
}

type DatabaseTypeVersionShowResponse struct {
	DatabaseTypeVersion DatabaseTypeVersion `json:"database_type_version"`
}

func (c Client) DatabaseTypeVersion(ctx context.Context, appID, addonID, versionID string) (DatabaseTypeVersion, error) {
	var res DatabaseTypeVersionShowResponse
	err := c.DBAPI(appID, addonID).DoRequest(ctx, &http.APIRequest{
		Method:   "GET",
		Endpoint: "/database_type_versions/" + versionID,
		Expected: http.Statuses{200},
	}, &res)
	if err != nil {
		return res.DatabaseTypeVersion, errgo.Notef(err, "fail to get database type version %v", versionID)
	}
	return res.DatabaseTypeVersion, nil
}
