package auth

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/go-scalingo"
)

const (
	ConfigVersionV2 = "2.0"
)

type ConfigData struct {
	AuthDataVersion   string          `json:"auth_data_version"`
	LastUpdate        time.Time       `json:"last_update"`
	AuthConfigPerHost json.RawMessage `json:"auth_config_data"`
}

func NewConfigData() *ConfigData {
	return &ConfigData{
		AuthDataVersion: ConfigVersionV2,
	}
}

type CredentialsData struct {
	TokenGenerator *scalingo.APITokenGenerator `json:"tokens"`
	User           *scalingo.User              `json:"user"`
}

type ConfigPerHostV2 map[string]*CredentialsData
