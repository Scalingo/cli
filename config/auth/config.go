package auth

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
)

const (
	ConfigVersionV1 = "1.0"
)

type ConfigData struct {
	AuthDataVersion   string          `json:"auth_data_version"`
	LastUpdate        time.Time       `json:"last_update"`
	AuthConfigPerHost json.RawMessage `json:"auth_config_data"`
}

func NewConfigData() *ConfigData {
	return &ConfigData{
		AuthDataVersion: ConfigVersionV1,
	}
}

type LegacyAuthConfigPerHost map[string]*LegacyUser

type LegacyUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	AuthToken string `json:"auth_token"`
}

type ConfigPerHostV1 map[string]*scalingo.User

func (config *ConfigData) MigrateToV1() error {
	var (
		legacyConfig LegacyAuthConfigPerHost
		newConfig    ConfigPerHostV1 = make(ConfigPerHostV1)
	)
	err := json.Unmarshal(config.AuthConfigPerHost, &legacyConfig)
	if err != nil {
		return errgo.Mask(err)
	}

	for host, legacyUser := range legacyConfig {
		newUser := &scalingo.User{
			ID: legacyUser.ID, Username: legacyUser.Username, FirstName: legacyUser.FirstName,
			LastName: legacyUser.LastName, Email: legacyUser.Email, AuthenticationToken: legacyUser.AuthToken,
		}
		newConfig[host] = newUser
	}

	config.AuthDataVersion = ConfigVersionV1
	c, err := json.Marshal(&newConfig)
	if err != nil {
		return errgo.Mask(err)
	}
	config.AuthConfigPerHost = json.RawMessage(c)
	config.LastUpdate = time.Now()

	return nil
}
