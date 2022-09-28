package auth

import (
	"encoding/json"
	"time"

	"github.com/Scalingo/go-scalingo/v6"
)

const (
	ConfigVersionV2  = "2.0"
	ConfigVersionV21 = "2.1"
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
	Tokens *UserToken     `json:"tokens"`
	User   *scalingo.User `json:"user"`
}

type ConfigPerHostV2 map[string]CredentialsData

type UserToken struct {
	Token string `json:"token"`
}
