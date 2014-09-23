package apps

import "time"

type App struct {
	Id    string `json:"_id"`
	Name  string `json:"name"`
	Owner struct {
		Email string `json:"email"`
		Id    string `json:"_id"`
	} `json:"owner"`
	GitUrl    string    `json:"git_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	LogsURL   string    `json:"logs_url"`
}

func (app App) String() string {
	return app.Name
}
