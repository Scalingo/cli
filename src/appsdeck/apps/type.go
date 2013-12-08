package apps

import (
	"fmt"
	"time"
)

type App struct {
	Id       string `json:"_id"`
	FullName string `json:"fullname"`
	Name     string `json:"name"`
	Owner    struct {
		Email string `json:"email"`
		Id    string `json:"_id"`
	} `json: "owner"`
	GitUrl    string    `json:"git_url"`
	CreatedAt time.Time `json: "created_at"`
	UpdatedAt time.Time `json: "update_at"`
}

func (app App) String() string {
	return fmt.Sprintf("%s \"%s\"", app.FullName, app.Name)
}
