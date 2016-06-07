package gitremote

import "fmt"

type ConfigFileError struct {
	Err error
}

func (err *ConfigFileError) Error() string {
	return fmt.Sprintf("GIT config file is not reachable: %v", err.Err)
}
