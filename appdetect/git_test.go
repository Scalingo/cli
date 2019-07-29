package appdetect

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestGetAppNameFromGitRemote(t *testing.T) {
	tests := map[string]struct {
		url             string
		expectedAppName string
	}{
		"Given a Git remote like on agora-fr1": {
			url:             "git@scalingo.com:my-app.git",
			expectedAppName: "my-app",
		},
		"Given a Git remote if SSH on a cutom port": {
			url:             "ssh://git@ssh.osc-fr1.scalingo.com:2200/my-app.git",
			expectedAppName: "my-app",
		},
		"Given a Git remote like on GitHub": {
			url:             "git@github.com:my-owner/my-app.git",
			expectedAppName: "my-app",
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			appName := getAppNameFromGitRemote(test.url)
			assert.Equal(t, test.expectedAppName, appName)
		})
	}
}
