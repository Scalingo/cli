package apps

import (
	"fmt"
)

type App struct {
	FullName string `json:"fullname"`
	Name     string `json:"name"`
}

func (app App) String() string {
	return fmt.Sprintf("%s (%s)", app.FullName, app.Name)
}
