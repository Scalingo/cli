package detect

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
)

func TestGetRegionFromGitRemote(t *testing.T) {
	tests := map[string]struct {
		url            string
		expectedRegion string
	}{
		"Given a Git remote like on agora-fr1": {
			url:            "ssh://git@ssh.osc-fr1.scalingo.com:22/my-app.git",
			expectedRegion: "osc-fr1",
		},
		"Given a Git remote if SSH on a cutom port": {
			url:            "ssh://git@ssh.osc-secnum-fr1.scalingo.com:22/my-app.git",
			expectedRegion: "osc-secnum-fr1",
		},
		"Given a Git remote like on GitHub": {
			url:            "git@github.com:my-owner/my-app.git",
			expectedRegion: "",
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			appName, _ := extractRegionFromGitRemote(test.url, &config.RegionsCache{
				Regions: []scalingo.Region{
					{
						Name: "osc-fr1",
						SSH:  "ssh.osc-fr1.scalingo.com:22",
					},
					{
						Name: "osc-secnum-fr1",
						SSH:  "ssh.osc-secnum-fr1.scalingo.com",
					},
				},
			})
			assert.Equal(t, test.expectedRegion, appName)
		})
	}
}
